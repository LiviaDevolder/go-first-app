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

func NewStoragePlayerInMemory() *StoragePlayerInMemory {
	return &StoragePlayerInMemory{map[string]int{}}
}
