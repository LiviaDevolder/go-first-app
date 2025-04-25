package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SketchPlayerStorage struct {
	scores     map[string]int
	winRecords []string
}

func (s *SketchPlayerStorage) GetPlayersPoints(name string) int {
	score := s.scores[name]
	return score
}

func (s *SketchPlayerStorage) RecordWin(name string) {
	s.winRecords = append(s.winRecords, name)
}

func TestGetPlayers(t *testing.T) {
	storage := SketchPlayerStorage{
		map[string]int{
			"Mary":  20,
			"Peter": 10,
		},
		nil,
	}
	server := &PlayerServer{&storage}

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
	storage := SketchPlayerStorage{
		map[string]int{},
		nil,
	}

	server := &PlayerServer{&storage}

	t.Run("record wins in HTTP POST request", func(t *testing.T) {
		// Arrange
		player := "Mary"

		request := newRequestPostRecordWin(player)
		response := httptest.NewRecorder()

		// Act
		server.ServeHTTP(response, request)

		// Assert
		verifyStatusCode(t, response.Code, http.StatusAccepted)

		if len(storage.winRecords) != 1 {
			t.Errorf("%d win records, expect %d", len(storage.winRecords), 1)
		}

		if storage.winRecords[0] != player {
			t.Errorf("didnt record correct player, result %s, expect %s", storage.winRecords[0], player)
		}
	})
}

func TestRecordWinsAndGetPoints(t *testing.T) {
	// Arrange
	storage := NewStoragePlayerInMemory()
	server := PlayerServer{storage}
	player := "Mary"

	// Act
	for range 3 {
		server.ServeHTTP(httptest.NewRecorder(), newRequestPostRecordWin(player))
	}

	response := httptest.NewRecorder()

	server.ServeHTTP(response, newRequestGetPoints(player))

	// Assert
	verifyStatusCode(t, response.Code, http.StatusOK)

	verifyRequestBody(t, response.Body.String(), "3")
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
