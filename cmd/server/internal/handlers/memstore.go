package handlers

import (
	"errors"
	"fmt"
	"sync"
)

type inMemStore struct {
	mu              *sync.RWMutex
	idCounter       uint64
	buffs           map[uint64]*Buff
	streamIDCounter uint64
	streams         map[uint64]Stream
}

//Those methods are here just for the successful code compilation
func (i *inMemStore) GetStream(id uint64) (Stream, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	s, ok := i.streams[id]
	if !ok {
		return Stream{}, fmt.Errorf("inMemStore#GetStream: %w", ErrNotFound)
	}
	return s, nil
}

func (i *inMemStore) SetStream(s Stream) (id uint64, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.streamIDCounter++
	s.ID = i.streamIDCounter
	i.streams[i.streamIDCounter] = s
	return i.streamIDCounter, nil
}

func (i *inMemStore) GetBuff(id uint64) (*Buff, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	b, ok := i.buffs[id]
	if !ok {
		return nil, errors.New("buff not found")
	}
	return b, nil
}

func (i *inMemStore) SetBuff(b *Buff) (uint64, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.idCounter++
	b.ID = i.idCounter
	i.buffs[i.idCounter] = b
	return i.idCounter, nil
}

func NewInMemStore() Store {
	return &inMemStore{
		mu:        &sync.RWMutex{},
		idCounter: 0,
		buffs:     make(map[uint64]*Buff),
		streams:   make(map[uint64]Stream),
	}
}
