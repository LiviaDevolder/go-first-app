package poker

import (
	"testing"
)

func TestPlayerFileStorageSystem(t *testing.T) {
	t.Run("/league reader", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Ari", "Wins": 10},
            {"Name": "Livia", "Wins": 0}]`)

		defer cleanDatabase()

		storage, err := NewPlayerFileStorageSystem(database)

		VerifyNoError(t, err)

		expect := []Player{
			{"Ari", 10},
			{"Livia", 0},
		}

		// Act
		result := storage.GetLeague()

		// Assert
		VerifyLeague(t, result, expect)
	})

	t.Run("get player score", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Ari", "Wins": 10},
            {"Name": "Livia", "Wins": 0}]`)

		defer cleanDatabase()

		storage, err := NewPlayerFileStorageSystem(database)

		VerifyNoError(t, err)

		expect := 0

		// Act
		result := storage.GetPlayerScore("Livia")

		// Assert
		VerifyScore(t, result, expect)
	})

	t.Run("store a player wins", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Livia", "Wins": 0},
            {"Name": "Ari", "Wins": 10}]`)

		defer cleanDatabase()

		expect := 1

		storage, err := NewPlayerFileStorageSystem(database)

		VerifyNoError(t, err)

		storage.SaveVictory("Livia")

		// Act
		result := storage.GetPlayerScore("Livia")

		// Assert
		VerifyScore(t, result, expect)
	})

	t.Run("store a new player win", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Livia", "Wins": 0},
            {"Name": "Ari", "Wins": 10}]`)

		defer cleanDatabase()

		expect := 1

		storage, err := NewPlayerFileStorageSystem(database)

		VerifyNoError(t, err)

		storage.SaveVictory("Anthony")

		// Act
		result := storage.GetPlayerScore("Anthony")

		// Assert
		VerifyScore(t, result, expect)
	})

	t.Run("deal with an empty file", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, "")
		defer cleanDatabase()

		// Act
		_, err := NewPlayerFileStorageSystem(database)

		// Assert
		VerifyNoError(t, err)
	})

	t.Run("order league", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Livia", "Wins": 0},
            {"Name": "Ari", "Wins": 10}]`)

		defer cleanDatabase()

		storage, _ := NewPlayerFileStorageSystem(database)

		expect := []Player{
			{"Ari", 10},
			{"Livia", 0},
		}

		// Act
		result := storage.GetLeague()

		// Assert
		VerifyLeague(t, result, expect)

		result = storage.GetLeague()
		VerifyLeague(t, result, expect)
	})
}

