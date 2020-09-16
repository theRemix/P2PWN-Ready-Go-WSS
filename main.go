package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/localtunnel/go-localtunnel"
)

// CHANGE ME
const appName = "p2pwn-ready-go-wss"
const displayName = "P2PWN Ready Go WSS"
const appRelease = "DEVELOPMENT"

// App  Config
// options to be sent to PTPWN Service during /api/connect
type appConfig struct {
	AppName     string `json:"app_name"`     // for grouping rooms in P2PWN
	DisplayName string `json:"display_name"` // used to display in P2PWN lobby
	Release     string `json:"release"`      // "PRODUCTION", "DEVELOPMENT"
	EntryURL    string `json:"entry_url"`    // url used as the entrypoint for your app, supplied by localtunnel
	Port        string // Server Listening Port
	P2pwnAddr   string // P2PWN Service Address
	// optional
	// HealthCheckURL string `json:"healthcheck_url"` // default: entry_url, if server cannot reach this endpoint, it will be unlisted
}

var Config = &appConfig{}

// P2PWN Service Config
// values returned from P2PWN Service response to /api/connect
type p2pwnConfig struct { // all values will be provided by P2PWN
	ID          string `json:"id"`           // public id assigned by P2PWN service
	AccessToken string `json:"access_token"` // private access token needed to perform actions on this host
	AppName     string `json:"app_name"`     // for grouping rooms in P2PWN
	DisplayName string `json:"display_name"` // used to display in P2PWN lobby
	EntryURL    string `json:"entry_url"`    // url used as the entrypoint for your app, supplied by localtunnel
}

var P2pwn = &p2pwnConfig{}

func main() {

	//================== P2PWN Setup Begin =====================//

	setConfig(&Config.AppName, "name", appName, "Name of this app")
	setConfig(&Config.Port, "port", "3000", "Port for server to listen on")
	setConfig(&Config.P2pwnAddr, "p2pwn", "https://p2pwn-production.herokuapp.com", "P2PWN Service Address")
	Config.DisplayName = displayName
	Config.Release = appRelease

	flag.Parse()

	port, portErr := strconv.Atoi(Config.Port)
	if portErr != nil {
		fmt.Printf("Invalid Port config: %s -> %v \n", Config.Port, port)
		os.Exit(1)
		return
	}

	lt, ltErr := localtunnel.Listen(localtunnel.Options{
		Subdomain: Config.AppName,
	})
	if ltErr != nil {
		fmt.Printf("Error creating localtunnel: %v\n", ltErr)
		os.Exit(1)
		return
	}

	Config.EntryURL = lt.URL()
	payload, _ := json.Marshal(Config)

	p2pwnRes, p2pwnErr := http.Post(Config.P2pwnAddr+"/api/connect", "application/json", bytes.NewBuffer(payload))
	if p2pwnErr != nil {
		fmt.Printf("Error Connecting to P2PWN Service: %v\n", p2pwnErr)
		os.Exit(1)
		return
	}

	defer p2pwnRes.Body.Close()
	if err := json.NewDecoder(p2pwnRes.Body).Decode(P2pwn); err != nil {
		fmt.Println("Unmarshal P2PWN Response Error:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("P2PWN is Ready: %+v\n", P2pwn)

	//------------------ P2PWN Setup Complete ---------------------//

	//================== Begin your app code  =====================//

	// Setup your handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello P2PWN-Go")
	})

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", server.Addr)
	server.Serve(lt)

	//------------------ End your app code -------------------------//
}

func setConfig(configPtr *string, flagName string, defaultVal string, help string) {
	flag.StringVar(configPtr, flagName, defaultVal, help)

	if val, ok := os.LookupEnv(flagName); ok {
		*configPtr = val
	}
}
