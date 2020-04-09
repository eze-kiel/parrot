package parrot

import (
	"bufio"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Client represents a client
type client struct {
	Nick string
	Conn net.Conn
}

// Message represents a message
type message struct {
	Sender  string
	Message string
}

type Server struct {
	Addr         string
	announcement chan client
	message      chan message
}

func (s *Server) Run() {
	s.announcement = make(chan client)
	s.message = make(chan message)

	l, err := net.Listen("tcp", s.Addr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", s.Addr)

	// close the socket when server is ended
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.orchestrator()
		go s.handleRequest(conn)
	}
}

func (s *Server) orchestrator() {
	var clients []net.Conn
	for {
		// Make the reception of annoucement non-blocking
		select {

		case newArrival := <-s.announcement:
			clients = append(clients, newArrival.Conn)
			log.Infof("New client: %s at %s\n", newArrival.Nick, newArrival.Conn.RemoteAddr().String())

			// Send to all clients that a new one arrived
			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("<server> " + newArrival.Nick + " (" + newArrival.Conn.RemoteAddr().String() + ") has joined the room\n")
				writer.Flush()
			}

		case newMessage := <-s.message:
			// Send to all clients the new message
			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("<" + newMessage.Sender + ">" + " " + newMessage.Message + "\n")
				writer.Flush()
			}
		}
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	var client client
	client.Conn = conn
	reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(client.Conn)
	// writer.Flush()
	// writer.WriteString("<server> enter your nickname and press ENTER: ")
	// writer.Flush()
	nick, _ := reader.ReadString('\n')

	client.Nick = strings.TrimSpace(nick)

	// Send new client to orchestrator
	s.announcement <- client

	close := func() {
		log.Print("closing connection")
		conn.Close()
	}

	defer close()

	for {
		var newMessage message

		newMessage.Sender = client.Nick

		msg, err := reader.ReadString('\n')

		if err == io.EOF {
			log.Print("client closed connection")
			break
		} else if err != nil {
			log.Fatal(err)
		}
		msg = strings.TrimSpace(msg)

		newMessage.Message = msg
		s.message <- newMessage
	}
}
