package main

import (
	"log"

	"github.com/namsral/flag"

	"github.com/eze-kiel/parrot/client"
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

	if serv == false && addr == "" {
		log.Fatal("Error: you must supply a server address")
	} else {
		addr = "127.0.0.1:3333"
	}

	if serv {

		server := &parrot.Server{
			Addr: addr,
		}

		server.Run()
	} else {
		client := &client.Client{
			Server: addr,
		}
		client.Run()
	}
}
