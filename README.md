# Parrot :bird: :speech_balloon:
![Licence](https://img.shields.io/badge/License-GPL-brightgreen)

A TUI irc-like chat written in go.

## Launch
```
options:
    -server
        default: false
        launch parrot as a server

    -addr <ip:port>
        default: 127.0.0.1:3333
        server mode: address of the server
        client mode: address to connect to
    
    -nick <nickname>
        nickname used to connect to a server
```

## Usage
Typical use as a server:

`parrot -server -addr 192.168.0.20:3333`

as a client:

`parrot -addr 192.168.0.20:3333 -nick JohnDoe`


## Roadmap
1. Add the following commands:
    * `/nick [new nick]` - change the nickname
    * `/who` - list all users
    * `/quit` - quit the room

## Author
Written by ezekiel.
