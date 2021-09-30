package LocalTypes

import (
	"container/list"
	"sync"
	"time"
)

type Request struct {
	Value string `json:"-"`
}

//Book type definitions
type Book struct {
	Value string `json:"-"`
	Owner string `json:"owner"`
	Writes int `json:"writes,omitempty"`
	Reads int `json:"reads,omitempty"`
	Age time.Time `json:"age,omitempty"`
}

//KvStore type definition. Mutex for lock/unlock when making operations on object
type KvStore struct {
	Depth int
	Books map[string]Book
	Order *list.List

	*sync.RWMutex
}
