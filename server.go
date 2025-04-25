package main

import (
	"fmt"
	"net/http"
)

type PlayerStorage interface {
	GetPlayersPoints(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	storage PlayerStorage
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.recordWin(w, player)
	case http.MethodGet:
		p.getScore(w, player)
	}
}

func (p *PlayerServer) getScore(w http.ResponseWriter, player string) {
	score := p.storage.GetPlayersPoints(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) recordWin(w http.ResponseWriter, player string) {
	p.storage.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
