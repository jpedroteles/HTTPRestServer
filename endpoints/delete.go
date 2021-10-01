package endpoints

import (
	"Week2Proj/constants"
	"Week2Proj/logger"
	"net/http"
	"strings"
)

//Delete given an isbn in path looks for it and if exists deletes it
func Delete(writer http.ResponseWriter, request *http.Request, auth string, s *StoreHandler) {
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	key := strings.TrimPrefix(request.URL.Path, constants.StorePath)
	key = strings.TrimLeft(key, "/")
	s.Store.Lock()
	book, ok := s.Store.Books[key]
	s.Store.Unlock()
	if !ok {
		logger.AppErrorLogger.Println("Book not found")
		msg := "404 key not found"
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	if book.Owner == auth || auth == constants.Admin {
		s.Store.Lock()
		//Remove key
		delete(s.Store.Books, key)
		s.Store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	} else {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}
}
