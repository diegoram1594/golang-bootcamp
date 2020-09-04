package main

import (
	"golangbootcamp/src/data"
	"golangbootcamp/src/server"
)

func main() {
	db := &data.DB{}
	db.NewDB()
	defer db.CloseDB()
	servidor := server.NewServer(*db)
	err := servidor.Listen()
	if err!= nil {
		panic(err)
	}



}






