package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/LiviaDevolder/go-first-app"
)

const fileNameDB = "game.db.json"

func main() {
	storage, close, err := poker.PlayerFileStorageSystemFromFile(fileNameDB)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	poker.NewCLI(storage, os.Stdin).PlayPoker()
}
