package endpoints

import (
	"Week2Proj/constants"
	"Week2Proj/logger"
	"net/http"
	"strings"
	"time"
)

//Get given an isbn in path look for it and return full object
func Get(writer http.ResponseWriter, request *http.Request, auth string, s *StoreHandler) {
	logger.AppInfoLogger.Println("Getting entry")
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	key := strings.TrimPrefix(request.URL.Path, constants.StorePath)
	key = strings.TrimLeft(key, "/")
	s.Store.Lock()
	book, ok := s.Store.Books[key]
	s.Store.Unlock()
	if !ok {
		logger.AppErrorLogger.Println("Book not found")
		msg := http.StatusText(http.StatusNotFound)
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	if book.Owner == auth {
		s.Store.Lock()
		book.Reads++
		book.Age = time.Now()
		s.Store.Books[key] = book
		s.Store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(book.Value))
	} else {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte(http.StatusText(http.StatusForbidden)))
	}

}
