package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const contentTypeJSON = "application/json"

type SketchPlayerStorage struct {
	scores     map[string]int
	winRecords []string
	league     League
}

func (s *SketchPlayerStorage) GetLeague() League {
	return s.league
}

func (s *SketchPlayerStorage) GetPlayersPoints(name string) int {
	score := s.scores[name]
	return score
}

func (s *SketchPlayerStorage) RecordWin(name string) {
	s.winRecords = append(s.winRecords, name)
}

func TestGetPlayers(t *testing.T) {
	database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Mary", "Wins": 20},
            {"Name": "Peter", "Wins": 10}]`)
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	verifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("get Mary's result", func(t *testing.T) {
		// Arrange
		request := newRequestGetPoints("Mary")
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusOK)
		verifyRequestBody(t, response.Body.String(), "20")
	})

	t.Run("get Peter's result", func(t *testing.T) {
		// Arrange
		request := newRequestGetPoints("Peter")
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusOK)
		verifyRequestBody(t, response.Body.String(), "10")
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
	database, cleanDatabase := createTmpFile(t, "[]")
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	verifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("record wins in HTTP POST request", func(t *testing.T) {
		// Arrange
		player := "Mary"

		request := newRequestPostRecordWin(player)
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusAccepted)

		if storage.GetPlayerScore(player) != 1 {
			t.Errorf("%d win records, expect %d", storage.GetPlayerScore(player), 1)
		}
	})
}

func TestRecordWinsAndGetPoints(t *testing.T) {
	// Arrange
	database, cleanDatabase := createTmpFile(t, "[]")
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)
	
	verifyNoError(t, err)
	
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
		verifyStatusCode(t, response.Code, http.StatusOK)

		verifyRequestBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		// Arrange
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		expect := []Player{
			{"Mary", 3},
		}

		// Act
		result := getLeagueAnswer(t, response.Body)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusOK)
		verifyLeague(t, result, expect)
	})
}

func TestLeague(t *testing.T) {
	database, cleanDatabase := createTmpFile(t, `[
            {"Name": "Livia", "Wins": 22},
						{"Name": "Ari", "Wins": 23},
            {"Name": "Anthony", "Wins": 22}]`)
	defer cleanDatabase()
	storage, err := NewPlayerFileStorageSystem(database)

	verifyNoError(t, err)

	server := NewPlayerServer(storage)

	t.Run("return 200 in /league", func(t *testing.T) {
		// Arrange
		request := newLeagueRequest()
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		result := getLeagueAnswer(t, response.Body)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusOK)
		verifyLeague(t, result, storage.GetLeague())
		verifyContentType(t, response, contentTypeJSON)
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

func verifyRequestBody(t *testing.T, result, expect string) {
	t.Helper()
	if result != expect {
		t.Errorf("result %s, expect %s", result, expect)
	}
}

func verifyStatusCode(t *testing.T, result, expect int) {
	t.Helper()
	if result != expect {
		t.Errorf("wrong status code, result %d, expected %d", result, expect)
	}
}

func getLeagueAnswer(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Cant parse %s player server answer: %v", body, err)
	}

	return
}

func verifyLeague(t *testing.T, result, expected []Player) {
	t.Helper()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("obtido %v esperado %v", result, expected)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func verifyContentType(t *testing.T, response *httptest.ResponseRecorder, expect string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != expect {
		t.Errorf("wrong response type, expect %s and got %s", expect, response.Result().Header)
	}
}
