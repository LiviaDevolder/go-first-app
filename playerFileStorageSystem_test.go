package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPlayerFileStorageSystem(t *testing.T) {
	t.Run("/league reader", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Ari", "Wins": 10},
            {"Name": "Livia", "Wins": 0}]`)

		defer cleanDatabase()

		storage, err := NewPlayerFileStorageSystem(database)

		verifyNoError(t, err)

		expect := []Player{
			{"Ari", 10},
			{"Livia", 0},
		}

		// Act
		result := storage.GetLeague()

		// Assert
		verifyLeague(t, result, expect)
	})

	t.Run("get player score", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Ari", "Wins": 10},
            {"Name": "Livia", "Wins": 0}]`)

		defer cleanDatabase()

		storage, err := NewPlayerFileStorageSystem(database)

		verifyNoError(t, err)

		expect := 0

		// Act
		result := storage.GetPlayerScore("Livia")

		// Assert
		verifyScore(t, result, expect)
	})

	t.Run("store a player wins", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Livia", "Wins": 0},
            {"Name": "Ari", "Wins": 10}]`)

		defer cleanDatabase()

		expect := 1

		storage, err := NewPlayerFileStorageSystem(database)

		verifyNoError(t, err)

		storage.SaveVictory("Livia")

		// Act
		result := storage.GetPlayerScore("Livia")

		// Assert
		verifyScore(t, result, expect)
	})

	t.Run("store a new player win", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Livia", "Wins": 0},
            {"Name": "Ari", "Wins": 10}]`)

		defer cleanDatabase()

		expect := 1

		storage, err := NewPlayerFileStorageSystem(database)

		verifyNoError(t, err)

		storage.SaveVictory("Anthony")

		// Act
		result := storage.GetPlayerScore("Anthony")

		// Assert
		verifyScore(t, result, expect)
	})

	t.Run("deal with an empty file", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, "")
		defer cleanDatabase()

		// Act
		_, err := NewPlayerFileStorageSystem(database)

		// Assert
		verifyNoError(t, err)
	})

	t.Run("order league", func(t *testing.T) {
		// Arrange
		database, cleanDatabase := createTmpFile(t, `[
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
		verifyLeague(t, result, expect)

		result = storage.GetLeague()
		verifyLeague(t, result, expect)
	})
}

func verifyScore(t *testing.T, result, expect int) {
	t.Helper()

	if result != expect {
		t.Errorf("result %d expect %d", result, expect)
	}
}

func createTmpFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("can not write temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func verifyNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expect no error but got %v", err)
	}
}
