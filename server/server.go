package server
import (
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"
)

// Basic type for the server.
// Contains the channels needed to communicate.
type Server struct {
	name        string
	Done        chan bool
	connections map[string]*Socket
	Local       *Socket
}

type Message struct {
	Content string
	Ip      string
}

type Socket struct {
	Conn  *websocket.Conn
	Pipe  chan *Message
	addr  string
	Alive chan bool
}

// Creates a server listening on localhost port 1234
func NewServer() *Server {
	server := &Server{
		name: "My Server",
		Done: make(chan bool),
		connections: make(map[string]*Socket),
		Local: nil,
	}
	go func() {
		http.Handle("/", websocket.Handler(createHandler(server)))
		err := http.ListenAndServe(":1234", nil)
		if err != nil {
			server.Done <- true
			panic("ListenAndServe: " + err.Error())
		}
	}()
	return server;
}
func createHandler(server *Server) (func(*websocket.Conn)) {
	return func(conn *websocket.Conn) {
		socket := server.AddConnection(conn)
		go handleInbound(socket)
		go handleOutbound(socket)
		<-socket.Alive
	}
}

func handleInbound(socket *Socket) {
	var err error
	for {
		var message string
		if err = websocket.Message.Receive(socket.Conn, &message); err != nil {
			fmt.Println("Cannot receive: " + err.Error())
			break
		}
		socket.SendMessage(message)
	}
}


func handleOutbound(socket *Socket) {
	select {
	case message := <-socket.Pipe:
		fmt.Println(message)
	}
}

func (socket *Socket) SendMessage(content string) {
	message := &Message{
		Content: content,
		Ip: socket.addr,
	}
	socket.Pipe <- message
}

func (server *Server) AddConnection(ws *websocket.Conn) *Socket {
	socket := &Socket{
		Conn: ws,
		Pipe: make(chan *Message),
		addr: ws.Request().RemoteAddr,
		Alive: make(chan bool),
	}
	if IsLocalhost(ws.Request().RemoteAddr) {
		fmt.Println("Local connection from ", socket.addr)
		server.Local = socket
	} else {
		fmt.Println("Remote connection from ", socket.addr)
		server.connections[socket.addr] = socket
	}
	return socket
}




