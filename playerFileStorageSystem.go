package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type PlayerFileStorageSystem struct {
	database *json.Encoder
	league   League
}

func (p *PlayerFileStorageSystem) GetLeague() League {
	sort.Slice(p.league, func(i, j int) bool {
		return p.league[i].Wins > p.league[j].Wins
	})

	return p.league
}

func (p *PlayerFileStorageSystem) GetPlayerScore(name string) int {
	player := p.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (p *PlayerFileStorageSystem) SaveVictory(name string) {
	player := p.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		p.league = append(p.league, Player{name, 1})
	}

	p.database.Encode(p.league)
}

func NewPlayerFileStorageSystem(file *os.File) (*PlayerFileStorageSystem, error) {
	err := initializePlayerBDDFile(file)

	if err != nil {
		return nil, fmt.Errorf("error initializing player file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("can not load player storage")
	}

	return &PlayerFileStorageSystem{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func initializePlayerBDDFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("trouble when using file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func PlayerFileStorageSystemFromFile(path string) (*PlayerFileStorageSystem, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	storage, err := NewPlayerFileStorageSystem(db)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file storage system, %v", err)
	}

	return storage, closeFunc, nil
}
