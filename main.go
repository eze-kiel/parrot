package main

import (
	"log"

	"github.com/eze-kiel/parrot/client"
	"github.com/eze-kiel/parrot/command"

	"github.com/namsral/flag"

	"github.com/eze-kiel/parrot/parrot"
)

func main() {
	var serv bool
	var nickname, addr string

	flag.BoolVar(&serv, "server", false, "server mode")
	flag.StringVar(&nickname, "nick", "guest", "nickname")
	flag.StringVar(&addr, "addr", "127.0.0.1:3333", "ip address")
	flag.Parse()

	if serv {
		runServer(addr)
	}

	runClient(addr, nickname)
}

func runClient(addr string, nickname string) {
	if addr == "" {
		log.Fatal("Error: you must supply a server address")
	}

	client := &client.Client{
		Server: addr,
		Nick:   nickname,
	}
	client.Run()
}

func runServer(addr string) {
	if addr == "" {
		addr = "127.0.0.1:3333"
	}

	server := &parrot.Server{
		Addr: addr,
	}

	commands := []parrot.Command{
		command.DateCommand{},
	}

	log.Fatal(server.Run(commands...))
}
