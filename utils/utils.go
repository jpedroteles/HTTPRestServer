package utils

import (
	"sync"

	LocalTypes "Week2Proj/Types"
	_ "Week2Proj/Types"
)

func SetUpData() *LocalTypes.StoreHandler {
	book := &LocalTypes.StoreHandler{
		Store: &LocalTypes.KvStore{
			Books:map[string]LocalTypes.Book{
				"d":{"a",1,"Get Set Go!", "John Smith", "user_a"},
				"e":{"b",2,"Be a Go Getter", "David Byrne", "user_b"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	return book
}
