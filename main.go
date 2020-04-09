package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

//Client represents a client
// type Client struct {
// 	Conn net.Conn
// }

func main() {
	announcement := make(chan net.Conn)
	message := make(chan string)

	l, err := net.Listen("tcp", ":3333")

	if err != nil {
		log.Fatal(err)
	}

	log.Print("listening on port 3333")

	// close the socket when server is ended
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go orchestrator(announcement)
		go handleRequest(conn, message, announcement)
	}
}

func orchestrator(announcement chan net.Conn) {
	var clients []net.Conn
	for {
		newArrival := <-announcement
		clients = append(clients, newArrival)
		for i := range clients {
			writer := bufio.NewWriter(clients[i])
			writer.WriteString("A new client arrived!: " + clients[i].RemoteAddr().String() + "\n")
			writer.Flush()
			fmt.Printf("A new client arrived!: %s\n", newArrival.RemoteAddr().String())
		}
	}
}

func handleRequest(conn net.Conn, message chan string, announcement chan net.Conn) {
	// Send new client to orchestrator
	announcement <- conn

	close := func() {
		log.Print("closing connection")
		conn.Close()
	}

	defer close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Print("client closed connection")
			break
		} else if err != nil {
			log.Panic(err)
		}
		msg = strings.TrimSpace(msg)

		if msg == "/q" {
			log.Printf("received stop signal")
			close()
			break
		} else {
			writer.WriteString("> " + msg + "\n")
			writer.Flush()
		}
	}
}
