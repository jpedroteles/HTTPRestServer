package endpoints

import (
	"Week2Proj/Types"
	"Week2Proj/constants"
	_ "Week2Proj/constants"
	"Week2Proj/utils"
	"encoding/json"
	"net/http"
	"strings"
)


//List lists all books or given an isbn in path only the one that matches
// Not sure if this method really need auth
func List(writer http.ResponseWriter, request *http.Request, auth string, s *LocalTypes.StoreHandler){
	key := strings.TrimPrefix(request.URL.Path,constants.ListPath)
	key = strings.TrimLeft(key,"/")
	if key ==""{
		//list all books
		s.Store.RLock()
		books := make([]LocalTypes.ListInfo, 0, len(s.Store.Books))
		for _, v := range s.Store.Books {
			//if v.Owner == auth || auth == Admin{
			books = append(books, LocalTypes.ListInfo{
				v.Key,
				v.Owner,
			})
			//}
		}
		s.Store.RUnlock()
		jsonBytes, err := json.Marshal(books)
		if err != nil {
			internalServerError(writer, request)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonBytes)
	}else{
		s.Store.RLock()
		book, ok := s.Store.Books[key]
		s.Store.RUnlock()
		if !ok{
			utils.AppErrorLogger.Println("404 key not found")
			msg := "404 key not found"
			http.Error(writer, msg, http.StatusNotFound)
			return
		}
		//if book.Owner == auth || auth == Admin{
		bookList := LocalTypes.ListInfo{
			book.Key,
			book.Owner,
		}
		jsonData, err := json.Marshal(bookList)
		if err != nil {
			internalServerError(writer, request)
			return
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(jsonData))
		}
		/*}else{
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Forbidden"))
		}*/
	}
}