package main

import (
	"golangbootcamp/src/data"
	"golangbootcamp/src/server"
)

func main() {
	db := &data.DB{}
	db.NewDB()
	server.InitDB(db)
	servidor := server.NewServer()
	err := servidor.Listen()
	if err!= nil {
		panic(err)
	}
	defer db.CloseDB()


}






