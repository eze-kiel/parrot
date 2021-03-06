package parrot

import (
	"bufio"
	"io"
	"net"
	"strings"
	"time"

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
	departure    chan client
	commands     []Command
}

func (s *Server) Run(c ...Command) error {
	s.announcement = make(chan client, 50)
	s.message = make(chan message, 50)
	s.departure = make(chan client)
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
		var waitingConn []net.Conn
		conn, err := l.Accept()

		waitingConn = append(waitingConn, conn)

		if err != nil {
			return err
		}

		rate := time.Second / 10
		burstLimit := 100
		tick := time.NewTicker(rate)
		defer tick.Stop()
		throttle := make(chan time.Time, burstLimit)
		go func() {
			for t := range tick.C {
				select {
				case throttle <- t:
				default:
				}
			} // does not exit after tick.Stop()
		}()

		for _, req := range waitingConn {
			<-throttle
			go s.handleRequest(req)
		}

	}
}

func (s *Server) orchestrator() {
	// var clients []net.Conn
	clients := make(map[string]net.Conn)
	for {
		// Make the reception of annoucement non-blocking
		select {

		case newArrival := <-s.announcement:
			// clients = append(clients, newArrival.Conn)
			clients[newArrival.Nick] = newArrival.Conn
			log.Infof("New client: %s at %s\n", newArrival.Nick, newArrival.Conn.RemoteAddr().String())

			// Send to all clients that a new one arrived
			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("<server> " + newArrival.Nick + " (" + newArrival.Conn.RemoteAddr().String() + ") has joined the room\n")
				writer.Flush()
			}

		case newDeparture := <-s.departure:

			// Delete the client which is leaving
			delete(clients, newDeparture.Nick)

			for i := range clients {
				writer := bufio.NewWriter(clients[i])
				writer.WriteString("<server> " + newDeparture.Nick + " (" + newDeparture.Conn.RemoteAddr().String() + ") has left the room\n")
				writer.Flush()
			}

		case newMessage := <-s.message:

			// Rate limiter for incoming messages
			rate := time.Second / 10
			burstLimit := 10
			tick := time.NewTicker(rate)
			defer tick.Stop()
			throttle := make(chan time.Time, burstLimit)
			go func() {
				for t := range tick.C {
					select {
					case throttle <- t:
					default:
					}
				}
			}()
			<-throttle

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

	nick, _ := reader.ReadString('\n')
	client.Nick = strings.TrimSpace(nick)

	// Send new client to orchestrator
	s.announcement <- client

	close := func() {
		log.Print("closing connection")
		s.departure <- client
		conn.Close()
	}

	defer close()

	for {
		var newMessage message

		newMessage.Sender = client.Nick

		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Print("client closed connection")

			return
		} else if err != nil {
			log.Fatal(err)
		}
		msg = strings.TrimSpace(msg)

		newMessage.Message = s.runCommand(msg)
		s.message <- newMessage

	}
}

func (s *Server) runCommand(msg string) string {
	if len(msg) == 0 || msg[0] != '/' {
		return msg
	}

	for _, cmd := range s.commands {
		res, err := cmd.Execute(msg)
		if err != nil {
			continue
		}

		return res
	}

	return msg
}
