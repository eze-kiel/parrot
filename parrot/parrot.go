package parrot

import (
	"bufio"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Command interface must be implemented by supported commands
type Command interface {
	Execute(string) (string, error)
}

type client struct {
	Nick string
	Conn net.Conn
}

type message struct {
	Sender  string
	Message string
}

type Server struct {
	Addr         string
	announcement chan client
	message      chan message
	commands     []Command
}

func (s *Server) Run(c ...Command) error {
	s.announcement = make(chan client, 50)
	s.message = make(chan message, 50)
	s.commands = c

	l, err := net.Listen("tcp", s.Addr)

	if err != nil {
		return err
	}

	log.Printf("listening on %s", s.Addr)

	// close the socket when server is ended
	defer l.Close()

	go s.orchestrator()

	for {
		conn, err := l.Accept()

		if err != nil {
			return err
		}

		go s.handleRequest(conn)
	}

	return nil
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
			// Debug purposes
			log.Infof("message transmitted: %s", newMessage)
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
			msg = " "
		} else if err != nil {
			log.Fatal(err)
		}
		msg = strings.TrimSpace(msg)

		newMessage.Message = s.runCommand(msg)
		s.message <- newMessage
	}
}

func (s *Server) runCommand(msg string) string {
	if msg[0] == '/' {
		for _, cmd := range s.commands {
			res, err := cmd.Execute(msg)
			if err != nil {
				continue
			}

			return res
		}
	}

	return msg
}
