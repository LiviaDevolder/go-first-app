package main

import (
	"log"
	"net/http"
)



func main() {
	server := &PlayerServer{NewStoragePlayerInMemory()}

	if err := http.ListenAndServe(":5001", server); err != nil {
		log.Fatalf("cant connect to port 5001 %v", err)
	}
}
