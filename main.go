package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("can not open file %s %v", dbFileName, err)
	}

	storage, err := NewPlayerFileStorageSystem(db)

	if err != nil {
		log.Fatalf("can not load player file system")
	}

	server := NewPlayerServer(storage)

	if err := http.ListenAndServe(":5001", server); err != nil {
		log.Fatalf("cant connect to port 5001 %v", err)
	}
}
