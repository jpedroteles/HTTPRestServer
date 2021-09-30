package endpoints

import (
	LocalTypes "Week2Proj/Types"
	"Week2Proj/constants"
	"io/ioutil"
	"net/http"
	"strings"
)

//CreateOrUpdate creates entry if isbn doesn't already exist, if exists then updates entry
func CreateOrUpdate(writer http.ResponseWriter, request *http.Request, auth string, s *StoreHandler) {
	var b LocalTypes.Book
	key := strings.TrimPrefix(request.URL.Path, constants.StorePath)
	key = strings.TrimLeft(key, "/")
	value, _ := ioutil.ReadAll(request.Body)
	s.Store.Lock()
	toUpdate, ok := s.Store.Books[key]
	s.Store.Unlock()
	//ISBN doesn't exist so we create a new entry
	if !ok {
		//Creation part
		s.Store.Lock()
		b.Value = string(value)
		b.Owner = auth
		s.Store.Books[key] = b
		s.Store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	} else {
		//Isbn exists, so we update instead of creating a new one
		if toUpdate.Owner == auth || auth == constants.Admin {
			s.Store.Lock()
			updateModel := LocalTypes.Book{
				string(value),
				toUpdate.Owner,
			}
			s.Store.Books[key] = updateModel
			s.Store.Unlock()
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Ok"))
		} else {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Forbidden"))
		}

	}
}
