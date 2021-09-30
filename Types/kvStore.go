package LocalTypes

import "sync"

type Request struct {
	Value string `json:"-"`
}

//Book type definitions
type Book struct {
	Value string `json:"-"`
	Owner string `json:"owner"`
}

//kvStore type definition. Mutex for lock/unlock when making operations on object
type KvStore struct {
	Books map[string]Book
	*sync.RWMutex
}
