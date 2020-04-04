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

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	log.Printf("accepting new connection %v", conn.RemoteAddr())

	close := func() {
		log.Print("closing connection")
		conn.Close()
	}

	defer close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Print("client closed connection")
			break
		} else if err != nil {
			log.Panic(err)
		}

		str = strings.TrimSpace(str)

		if str == "/q" {
			log.Printf("received stop signal")
			close()
			break
		} else {
			writer.WriteString(fmt.Sprintf("> %s\n", str))
			writer.Flush()
		}
	}
}
