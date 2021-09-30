package server

import (
	LocalTypes "Week2Proj/Types"
	"Week2Proj/constants"
	"Week2Proj/endpoints"
	"net/http"
	"strings"
	"sync"
)

func SetUpServer() *http.ServeMux {
	mux := http.NewServeMux()
	data := SetUpData()
	mux.HandleFunc(constants.PingPath, endpoints.Ping)
	mux.HandleFunc(constants.ShutdownPath, endpoints.Shutdown)
	mux.Handle(constants.StorePath, data)
	mux.Handle(strings.TrimRight(constants.StorePath, "/"), data)
	mux.Handle(constants.ListPath, data)
	mux.Handle(constants.ListPath+"/", data)
	return mux
}

func SetUpData() *endpoints.StoreHandler {
	book := &endpoints.StoreHandler{
		Store: &LocalTypes.KvStore{
			Books:   map[string]LocalTypes.Book{},
			RWMutex: &sync.RWMutex{},
		},
	}
	return book
}
