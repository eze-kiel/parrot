package main

import (
	"bufio"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {

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

		// Crée une goroutine à chaque connexion
		go handleRequest(conn, message)
	}
}

func handleRequest(conn net.Conn, message chan string) {
	log.Printf("accepting new connection %v", conn.RemoteAddr())

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
		message <- msg
		msg = <-message

		if msg == "/q" {
			log.Printf("received stop signal")
			close()
			break
		} else {
			writer.WriteString(<-message)
			writer.Flush()
		}
	}
}
