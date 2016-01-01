package main
import (
	"fmt"
	"vote-server/server"
)

func main() {
	s := server.NewServer()
	fmt.Println("Started the server")

	<- s.Done
}
