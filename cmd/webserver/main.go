package main

import (
	"log"
	"net/http"
	poker "github.com/LiviaDevolder/go-first-app"
)

const fileNameDB = "game.db.json"

func main() {
	storage, close, err := poker.PlayerFileStorageSystemFromFile(fileNameDB)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(storage)

	if err := http.ListenAndServe(":5001", server); err != nil {
		log.Fatalf("cant connect to port 5001 %v", err)
	}
}
