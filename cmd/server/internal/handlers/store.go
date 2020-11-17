package handlers

import (
	"errors"
)

var ErrNotFound = errors.New("Entity not found")

type BuffStore interface {
	GetBuff(id uint64) (*Buff, error)
	SetBuff(*Buff) (id uint64, err error)
}

//Please notice, that store interface doesn't return pointer to Stream on purpose.
//Pointers exists to share state, and here we aren't share any state - we want value
//And also widely used pointers can result in more heap allocations and lead to
//GC overloading
type StreamStore interface {
	GetStream(id uint64) (Stream, error)
	SetStream(Stream) (id uint64,err  error)
}

type Store interface {
	BuffStore
	StreamStore
}
