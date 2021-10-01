package endpoints

import (
	"Week2Proj/Types"
	"Week2Proj/constants"
	"Week2Proj/logger"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

//List lists all books or given an isbn in path only the one that matches
// Not sure if this method really need auth
func List(writer http.ResponseWriter, request *http.Request, _ string, s *StoreHandler) {
	logger.AppInfoLogger.Println("List entries")
	writer.Header().Set("content-type", "application/json")
	key := strings.TrimPrefix(request.URL.Path, constants.ListPath)
	key = strings.TrimLeft(key, "/")
	if key == "" {
		//list all books
		s.Store.Lock()
		defer s.Store.Unlock()
		books := make([]LocalTypes.ListInfo, 0, len(s.Store.Books))
		for key, v := range s.Store.Books {
			books = append(books, LocalTypes.ListInfo{
				Key:    key,
				Owner:  v.Owner,
				Reads:  v.Reads,
				Writes: v.Writes,
				Age:    AgeMilli(v.Age),
			})
		}
		jsonBytes, err := json.Marshal(books)
		if err != nil {
			internalServerError(writer, request)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonBytes)
	} else {
		s.Store.Lock()
		book, ok := s.Store.Books[key]
		defer s.Store.Unlock()
		if !ok {
			logger.AppErrorLogger.Println("404 key not found")
			msg := "404 key not found"
			http.Error(writer, msg, http.StatusNotFound)
			return
		}
		bookList := LocalTypes.ListInfo{
			Key:    key,
			Owner:  book.Owner,
			Reads:  book.Reads,
			Writes: book.Writes,
			Age:    AgeMilli(book.Age),
		}
		jsonData, err := json.Marshal(bookList)
		if err != nil {
			internalServerError(writer, request)
			return
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write(jsonData)
		}
	}
}

func AgeMilli(age time.Time) int64 {
	test := time.Now().Sub(age)
	return test.Milliseconds()
}
