# Parrot :bird: :speech_balloon:
![Licence](https://img.shields.io/badge/License-GPL-brightgreen)

![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)

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

    -sound
        default: false
        play a sound a every new message
```

Note that when a client doesn't provide a nickname, a random number will be appended to 'guest' in order to avoid duplicates nicknames.

## Usage
Typical use as a server:

`parrot -server -addr 192.168.0.20:3333`

as a client:

`parrot -addr 192.168.0.20:3333 -nick JohnDoe -sound`

## Commands
```
    /date
        print the date following "Monday, 2006/01/02" format
```
## Roadmap :soon:
* [x] ~~Manage rate limiting~~
* [x] ~~Add notification sound~~
* [ ] Add the following commands:
    * `/nick [new nick]` - change the nickname
    * `/who` - list all users
    * `/quit` - quit the room
* [ ] Write tests (in progress...)
* [ ] Add encryption
* [ ] Embed notification sound in the binary file
* [ ] Add a disconnect message
* [ ] Fix issue#1

## Author
Written by ezekiel.
