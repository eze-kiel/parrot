package main

import (
	"bufio"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Client represents a client
type Client struct {
	Nick string
	Conn net.Conn
}

// Message represents a message
type Message struct {
	Sender  string
	Message string
}

func main() {
	announcement := make(chan Client)
	message := make(chan Message)

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

		go orchestrator(announcement, message)
		go handleRequest(conn, message, announcement)
	}
}

func orchestrator(announcement chan Client, message chan Message) {
	var clients []net.Conn
	for {
		// Make the reception of annoucement non-blocking
		select {

		case newArrival := <-announcement:
			clients = append(clients, newArrival.Conn)
			log.Infof("A new client arrived!: %s\n", newArrival.Conn.RemoteAddr().String())

			// Send to all clients that a new one arrived
			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("A new client arrived!: " + newArrival.Conn.RemoteAddr().String() + "\n")
				writer.Flush()
			}

		case newMessage := <-message:
			// Send to all clients the new message
			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("<" + newMessage.Sender + ">" + " " + newMessage.Message + "\n")
				writer.Flush()
			}
		}
	}
}

func handleRequest(conn net.Conn, message chan Message, announcement chan Client) {
	var client Client
	client.Conn = conn
	reader := bufio.NewReader(conn)
	nick, _ := reader.ReadString('\n')
	client.Nick = strings.TrimSpace(nick)

	// Send new client to orchestrator
	announcement <- client

	close := func() {
		log.Print("closing connection")
		conn.Close()
	}

	defer close()

	for {
		var newMessage Message

		newMessage.Sender = client.Nick

		msg, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Print("client closed connection")
			break
		} else if err != nil {
			log.Panic(err)
		}
		msg = strings.TrimSpace(msg)

		newMessage.Message = msg
		message <- newMessage

		if msg == "/q" {
			log.Printf("received stop signal")
			close()
			break
		} else {
			message <- newMessage
		}
	}
}
