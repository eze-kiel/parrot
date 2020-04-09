package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	announcement := make(chan string)
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

func orchestrator(announcement chan string) {
	newArrival := <-announcement
	fmt.Printf("A new client arrived!: %s\n", newArrival)
}

func handleRequest(conn net.Conn, message chan string, announcement chan string) {
	//log.Printf("accepting new connection %v", conn.RemoteAddr())
	// Send IP to orchestrator
	announcement <- conn.RemoteAddr().String()

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
