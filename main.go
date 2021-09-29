package main

import (
	"Week2Proj/server"
	_ "Week2Proj/server"
	"Week2Proj/utils"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main(){
	var(
		port int
	)
	//Handling arguments
	flag.IntVar(&port, "port", 8000, "port to listen on")

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			_,err:= strconv.Atoi(f.Value.String())
			if err != nil{
				utils.AppErrorLogger.Println("Failure to parse to int ", err)
				os.Exit(-1)
			}
		}
	})

	flag.Parse()

	mux := server.SetUpServer()

	utils.SetUpLogger()
	utils.AppInfoLogger.Println("Starting up proj in port: "+ strconv.Itoa(port))

	portString:= strconv.Itoa(port)
	err:= http.ListenAndServe(":"+portString, mux)
	if err != nil {
		utils.AppErrorLogger.Println("Failure to bind to the port ", err)
		os.Exit(-2)
	}
}

//Get given an isbn in path look for it and return full object
func (s storeHandler) Get(writer http.ResponseWriter, request *http.Request, auth string) {
	key := strings.TrimPrefix(request.URL.Path,StorePath)
	key = strings.TrimLeft(key,"/")
	s.store.RLock()
	book, ok := s.store.books[key]
	s.store.RUnlock()
	if !ok{
		utils.AppErrorLogger.Println("Book not found")
		msg := "404 key not found"
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	if book.Owner==auth {
		jsonData, err := json.Marshal(book)
		if err != nil {
			internalServerError(writer, request)
			return
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(jsonData))
		}
	}else{
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}

}

//Delete given an isbn in path looks for it and if exists deletes it
func (s storeHandler) Delete(writer http.ResponseWriter, request *http.Request, auth string) {
	key := strings.TrimPrefix(request.URL.Path,StorePath)
	key = strings.TrimLeft(key,"/")
	s.store.RLock()
	book, ok := s.store.books[key]
	s.store.RUnlock()
	if !ok{
		utils.AppErrorLogger.Println("Book not found")
		msg := "404 key not found"
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	if book.Owner == auth || auth == Admin {
		s.store.Lock()
		//Remove key
		delete(s.store.books, key)
		s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	}else{
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}

}

//CreateOrUpdate creates entry if isbn doesn't already exist, if exists then updates entry
func (s storeHandler) CreateOrUpdate(writer http.ResponseWriter, request *http.Request, auth string) {
	var b Book
	if err := json.NewDecoder(request.Body).Decode(&b); err != nil {
		internalServerError(writer, request)
		return
	}
	s.store.RLock()
	_, ok := s.store.books[string(b.ISBN)]
	s.store.RUnlock()
	//ISBN doesn't exist so we create a new entry
	if !ok{
		//Creation part
		s.store.Lock()
		b.Owner = auth
		s.store.books[strconv.Itoa(b.ISBN)] = b
		s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	}else{
		//Isbn exists, so we update instead of creating a new one
		if b.Owner == auth || auth == Admin{
			s.store.Lock()
			updateModel := Book{
				b.Key,
				b.ISBN,
				b.Title,
				b.Author,
				b.Owner,
			}
			s.store.books[string(b.ISBN)] = updateModel
			s.store.Unlock()
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Ok"))
		}else{
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Forbidden"))
		}

	}
}


