# P2PWN Ready Go WSS

P2PWN WSS Starter Kit for Go

Running this project will perform 2 tasks

## 1. Create a new Room in the P2PWN Lobby

named `P2PWN Ready Go WSS`

see all rooms here: https://p2pwn-production.herokuapp.com/lobby

In this lobby, the `P2PWN Ready Go WSS` will link to the `entry_url` _for example `https://p2pwn-ready-go-wss.loca.lt`_

## 2. Run a demo server

this demo server is listening on port `3000` on the P2PWN connected localtunnel interface.
you and everyone else in the world can access your demo server at the published `entry_url`

```sh
curl https://p2pwn-ready-go-wss.loca.lt/

Hello P2PWN-Go
```

## Dependencies

```sh
make env
make deps
```
edit `.env` if needed

## Run

```sh
make run
```
