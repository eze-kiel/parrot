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
        default: guest
        nickname used to connect to a server
```

Note that when a client doesn't provide a nickname, a random number will be appended to 'guest' in order to avoid duplicates nicknames.
## Usage
Typical use as a server:

`parrot -server -addr 192.168.0.20:3333`

as a client:

`parrot -addr 192.168.0.20:3333 -nick JohnDoe`

## Commands
```
    /date
        print the date following "Monday, 2006/01/02" format
```
## Roadmap
1. Append a random number to the 'guest' username to dissociate them
1. Add the following commands:
    * `/nick [new nick]` - change the nickname
    * `/who` - list all users
    * `/quit` - quit the room

## Author
Written by ezekiel.
