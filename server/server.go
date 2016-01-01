package server
import (
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"
)

// Basic type for the server.
// Contains the two channels needed to close the server and send/receive messages
type Server struct {
	Done        chan bool
	connections map[string]*socket
	Pipe        chan *Message
}

type socket struct {
	Conn  *websocket.Conn
	Host  string
	Alive chan bool
}

// Creates a server listening on localhost port 1234
func NewServer() *Server {
	server := &Server{
		Done: make(chan bool),
		connections: make(map[string]*socket),
		Pipe: make(chan *Message),
	}
	go server.startHttpListener()
	go server.handleOutbound()

	return server;
}

// Send the signal to close the server and all connections
func (server *Server) CloseConnections() {
	server.Done <- true
}



func createHandler(server *Server) (func(*websocket.Conn)) {
	return func(conn *websocket.Conn) {
		socket := server.addConnection(conn)
		server.Pipe <- NewMessage(socket.Host, "", CONNECTION_ESTABLISHED)
		go server.handleInbound(socket)
		<-socket.Alive
	}
}

func (server *Server) startHttpListener() {
	http.Handle("/", websocket.Handler(createHandler(server)))
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		server.Done <- true
		panic("ListenAndServe: " + err.Error())
	}
}

func (server *Server) handleInbound(socket *socket) {
	var err error
	for {
		var message string
		if err = websocket.Message.Receive(socket.Conn, &message); err != nil {
			fmt.Println("Cannot receive: " + err.Error())
			break
		}
		server.Pipe <- NewMessage(socket.Host, message, INBOUND)
	}
}


func (server *Server) handleOutbound() {
	for {
		select {
		case message := <-server.Pipe:
			if socket := server.connections[message.Host]; socket != nil {
				fmt.Println("Sending " + message.Content + " to " + socket.Host)
				fmt.Fprint(socket.Conn, message.Content)
			}
		case <-server.Done:
			fmt.Println("stopping outbound channel")
		}
	}
}


// Actually close the server and the connections
func (server *Server) onCloseConnections() {
	<-server.Done
	for _, socket := range server.connections {
		socket.Alive <- false
	}
}


func (server *Server) addConnection(ws *websocket.Conn) *socket {
	socket := &socket{
		Conn: ws,
		Host: ws.Request().RemoteAddr,
		Alive: make(chan bool),
	}
	server.connections[socket.Host] = socket
	fmt.Println("Connection from ", socket.Host)
	return socket
}
