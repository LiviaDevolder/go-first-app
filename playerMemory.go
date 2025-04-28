package main

type StoragePlayerInMemory struct {
	storage map[string]int
}

func (s *StoragePlayerInMemory) GetPlayersPoints(name string) int {
	return s.storage[name]

}

func (s *StoragePlayerInMemory) RecordWin(name string) {
	s.storage[name]++
}

func (s *StoragePlayerInMemory) GetLeague() []Player {
	var league []Player

	for name, wins := range s.storage {
		league = append(league, Player{name, wins})
	}

	return league
}

func NewStoragePlayerInMemory() *StoragePlayerInMemory {
	return &StoragePlayerInMemory{map[string]int{}}
}
