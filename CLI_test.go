package poker_test

import (
	"strings"
	"testing"

	poker "github.com/LiviaDevolder/go-first-app"
)

func TestCLI(t *testing.T) {
	// Arrange
	database, cleanDatabase := poker.CreateTmpFile(t, `[]`)
	defer cleanDatabase()

	in := strings.NewReader("Ari wins\n")
	playerStorage, _ := poker.NewPlayerFileStorageSystem(database)
	expect := "Ari"

	cli := poker.NewCLI(playerStorage, in)

	// Act
	cli.PlayPoker()

	// Assert
	poker.VerifyPlayerWin(t, playerStorage, expect)
}
