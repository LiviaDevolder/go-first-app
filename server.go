package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PlayerStorage interface {
	GetPlayersPoints(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	storage *PlayerFileStorageSystem
	http.Handler
}

type Player struct {
	Name string
	Wins int
}

func NewPlayerServer(storage *PlayerFileStorageSystem) *PlayerServer {
	p := new(PlayerServer)

	p.storage = storage

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(p.handleLeague))
	router.Handle("/players/", http.HandlerFunc(p.handlePlayers))

	p.Handler = router

	return p
}

func (p *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.storage.GetLeague())
}

func (p *PlayerServer) handlePlayers(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.recordWin(w, player)
	case http.MethodGet:
		p.getScore(w, player)
	}
}

func (p *PlayerServer) getScore(w http.ResponseWriter, player string) {
	score := p.storage.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) recordWin(w http.ResponseWriter, player string) {
	p.storage.SaveVictory(player)
	w.WriteHeader(http.StatusAccepted)
}
