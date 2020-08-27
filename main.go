package main

import (
	"fmt"
	"golangbootcamp/src/server"
)

func main() {
	servidor := server.NewServer()
	err := servidor.Listen()
	if err!= nil {
		fmt.Println(err)
	}
}






