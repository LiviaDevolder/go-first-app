package poker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const contentTypeJSON = "application/json"

func TestGetPlayers(t *testing.T) {
	database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Mary", "Wins": 20},
            {"Name": "Peter", "Wins": 10}]`)
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	VerifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("get Mary's result", func(t *testing.T) {
		// Arrange
		request := newRequestGetPoints("Mary")
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusOK)
		VerifyRequestBody(t, response.Body.String(), "20")
	})

	t.Run("get Peter's result", func(t *testing.T) {
		// Arrange
		request := newRequestGetPoints("Peter")
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusOK)
		VerifyRequestBody(t, response.Body.String(), "10")
	})

	t.Run("throws 404 when player doesnt exists", func(t *testing.T) {
		// Arrange
		request := newRequestGetPoints("George")
		response := httptest.NewRecorder()
		expect := http.StatusNotFound

		// Act
		server.ServeHTTP(response, request)
		result := response.Code

		// Assert
		if result != expect {
			t.Errorf("status %d but expected %d", result, expect)
		}
	})
}

func TestWinsStorage(t *testing.T) {
	database, cleanDatabase := CreateTmpFile(t, "[]")
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	VerifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("record wins in HTTP POST request", func(t *testing.T) {
		// Arrange
		player := "Mary"

		request := newRequestPostRecordWin(player)
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusAccepted)

		if storage.GetPlayerScore(player) != 1 {
			t.Errorf("%d win records, expect %d", storage.GetPlayerScore(player), 1)
		}
	})
}

func TestRecordWinsAndGetPoints(t *testing.T) {
	// Arrange
	database, cleanDatabase := CreateTmpFile(t, "[]")
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	VerifyNoError(t, err)

	server := NewPlayerServer(storage)
	player := "Mary"

	// Act
	for range 3 {
		server.ServeHTTP(httptest.NewRecorder(), newRequestPostRecordWin(player))
	}

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newRequestGetPoints(player))

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusOK)

		VerifyRequestBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		// Arrange
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())

		expect := []Player{
			{"Mary", 3},
		}

		// Act
		result := GetLeagueAnswer(t, response.Body)

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusOK)
		VerifyLeague(t, result, expect)
	})
}

func TestLeague(t *testing.T) {
	database, cleanDatabase := CreateTmpFile(t, `[
            {"Name": "Livia", "Wins": 22},
						{"Name": "Ari", "Wins": 23},
            {"Name": "Anthony", "Wins": 22}]`)
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	VerifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("return 200 in /league", func(t *testing.T) {
		// Arrange
		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		result := GetLeagueAnswer(t, response.Body)

		// Assert
		VerifyStatusCode(t, response.Code, http.StatusOK)
		VerifyLeague(t, result, storage.GetLeague())
		VerifyContentType(t, response, contentTypeJSON)
	})
}

func newRequestGetPoints(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newRequestPostRecordWin(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}
