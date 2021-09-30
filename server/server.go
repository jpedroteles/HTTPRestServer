package server

import (
	LocalTypes "Week2Proj/Types"
	"Week2Proj/constants"
	"Week2Proj/endpoints"
	"container/list"
	"net/http"
	"strings"
	"sync"
)

func SetUpServer(depth int) *http.ServeMux {
	mux := http.NewServeMux()
	data := SetUpData(depth)
	mux.HandleFunc(constants.PingPath, endpoints.Ping)
	mux.HandleFunc(constants.ShutdownPath, endpoints.Shutdown)
	mux.Handle(constants.StorePath, data)
	mux.Handle(strings.TrimRight(constants.StorePath, "/"), data)
	mux.Handle(constants.ListPath, data)
	mux.Handle(constants.ListPath+"/", data)
	return mux
}

func SetUpData(depth int) *endpoints.StoreHandler {
	if depth == 0{
		depth = constants.DefaultDepth
	}

	book := &endpoints.StoreHandler{
		Store: &LocalTypes.KvStore{
			Books:   map[string]LocalTypes.Book{},
			RWMutex: &sync.RWMutex{},
			Order: list.New(),
			Depth: depth,
		},
	}
	return book
}
