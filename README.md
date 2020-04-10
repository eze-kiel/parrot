# Parrot :bird: :speech_balloon:
![Licence](https://img.shields.io/badge/License-GPL-brightgreen)

A TUI irc-like chat written in go

## Launch
```Bash
Server mode :
    parrot -server
        server on default listen (127.0.0.1:3333)
    parrot -server 1.2.3.4:3333
        server listening on 1.2.3.4:3333

Client mode :
    parrot 1.2.3.4:3333
        client connection to 1.2.3.4:3333
```

## Usage
The first message you write will be your nickname. If you want to quit the TUI, just press `Esc`.

## Roadmap
1. Add the following commands:
    * `/nick [new nick]` - change the nickname
    * `/who` - list all users
    * `/quit` - quit the room

## Author
Written by ezekiel.
