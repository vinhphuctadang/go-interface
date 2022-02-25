package main

import (
	"errors"
	"log"
	"os"

	"github.com/vinhphuctadang/go-interface/db"
	mongoimpl "github.com/vinhphuctadang/go-interface/db/mongo_impl"
	mysqlimpl "github.com/vinhphuctadang/go-interface/db/mysql_impl"
)

func panicIf(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func main() {
	dbType := os.Getenv("DB_TYPE")
	dbUri := os.Getenv("DB_CONNECTION_URI")

	var dbSvc db.DBService
	var err error

	switch dbType {
	case "mongo":
		dbSvc, err = mongoimpl.NewDbServiceMongoBackend(dbUri)
	case "mysql":
		dbSvc, err = mysqlimpl.NewDbServiceMySqlBackend(dbUri, 10)
	default:
		err = errors.New("unsupported db")
	}

	panicIf(err)
	dbSvc.CreateAccount("hello", "word")
}
