package main
import (
	"fmt"
	"vote-server/server"
)

func main() {
	s := server.NewServer()
	fmt.Println("Started the server")

	go echo(s)

	<-s.Done
}

func echo(server *server.Server) {

	for {
		message := <-server.Pipe
		fmt.Println("Received ", message)
		server.Pipe <- message
	}
}
