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

	addr := ""

	flag.BoolVar(&serv, "server", false, "server mode")
	flag.Parse()

	if flag.NArg() > 0 {
		addr = flag.Arg(0)
	}

	if serv {
		runServer(addr)
	}

	runClient(addr)
}

func runClient(addr string) {
	if addr == "" {
		log.Fatal("Error: you must supply a server address")
	}

	client := &client.Client{
		Server: addr,
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
