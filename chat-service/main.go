package main

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/api"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/config"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/repository"
	"github.com/go-sql-driver/mysql"
	"log"
)

func main() {

	db, err := repository.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(config.Envs.Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to storage")
}
