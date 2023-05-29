package pg

import (
	"database/sql"
	"fmt"
	"github.com/gookit/ini/v2"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
)

var db *sql.DB
var logger *logrus.Entry

func Bootstrap() {
	logger = logrus.WithFields(logrus.Fields{})
	logger.Info("Connecting to PG")

	dbConfig := ini.StringMap("db")
	port, _ := strconv.Atoi(dbConfig["port"])
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"], port, dbConfig["user"], dbConfig["password"], dbConfig["name"], dbConfig["sslMode"])
	if d, err := sql.Open("postgres", psqlInfo); err != nil {
		logger.Fatal("Failed to connect: ", err)
	} else {
		db = d
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping: ", err)
	}
	log.Println("Connected")
}

func Exit() {
	_ = db.Close()
}
