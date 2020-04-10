package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eze-kiel/parrot/client"
	"github.com/eze-kiel/parrot/command"

	"github.com/namsral/flag"

	"github.com/eze-kiel/parrot/parrot"
)

func main() {
	var serv bool
	var nickname, addr string

	flag.BoolVar(&serv, "server", false, "server mode")
	flag.StringVar(&nickname, "nick", "", "nickname")
	flag.StringVar(&addr, "addr", "127.0.0.1:3333", "ip address")
	flag.Parse()

	if serv {
		runServer(addr)
	}

	if nickname == "" {
		fmt.Println("You must provide a nickname with the flag -nick")
		os.Exit(1)
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
