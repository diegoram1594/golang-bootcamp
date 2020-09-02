package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golangbootcamp/src/model"
)


var errDB error

type DB struct {
	database *sql.DB
}

type IDataBase interface {
	NewDB()
	CloseDB()
}

type ResultProductUser struct {
	ProductUser model.ProductUser
	Error error
}

func (db *DB) NewDB()  {
	db.database, errDB = sql.Open("mysql", "root:my-secret-pw@/bootcamp")
	if errDB != nil {
		panic(errDB.Error())
	}
}
func (db *DB) CloseDB()  {
	err := db.database.Close()
	if err != nil{
		panic(err)
	}
}





