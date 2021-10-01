package endpoints

import (
	LocalTypes "Week2Proj/Types"
	"Week2Proj/constants"
	"Week2Proj/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//CreateOrUpdate creates entry if isbn doesn't already exist, if exists then updates entry
func CreateOrUpdate(writer http.ResponseWriter, request *http.Request, auth string, s *StoreHandler) {
	logger.AppInfoLogger.Println("Putting entry")
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	key := strings.TrimPrefix(request.URL.Path, constants.StorePath)
	key = strings.TrimLeft(key, "/")
	value, _ := ioutil.ReadAll(request.Body)
	s.Store.Lock()
	toUpdate, ok := s.Store.Books[key]
	s.Store.Unlock()
	//ISBN doesn't exist so we create a new entry
	if !ok {
		//Creation part
		b := LocalTypes.Book{
			Value:  string(value),
			Owner:  auth,
			Writes: 1,
			Age:    time.Now(),
		}
		s.Store.Lock()
		s.Store.Books[key] = b
		s.Store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(http.StatusText(http.StatusOK)))
	} else {
		//Isbn exists, so we update instead of creating a new one
		if toUpdate.Owner == auth || auth == constants.Admin {
			s.Store.Lock()
			updateModel := LocalTypes.Book{
				Value:  string(value),
				Owner:  toUpdate.Owner,
				Writes: toUpdate.Writes + 1,
				Age:    time.Now(),
			}
			s.Store.Books[key] = updateModel
			s.Store.Unlock()
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(http.StatusText(http.StatusOK)))
		} else {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(http.StatusText(http.StatusForbidden)))
		}

	}
}
