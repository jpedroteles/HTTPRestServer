package LocalTypes

import "sync"

//Book type definitions
type Book struct {
	Key    string `json:"key"`
	ISBN   int    `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Owner  string `json:"owner"`
}

//kvStore type definition. Mutex for lock/unlock when making operations on object
type KvStore struct {
	Books map[string]Book
	*sync.RWMutex
}
