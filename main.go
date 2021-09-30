package main

import (
	"Week2Proj/server"
	_ "Week2Proj/server"
	"Week2Proj/utils"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		port int
	)
	//Handling arguments
	flag.IntVar(&port, "port", 8000, "port to listen on")

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			_, err := strconv.Atoi(f.Value.String())
			if err != nil {
				utils.AppErrorLogger.Println("Failure to parse to int ", err)
				os.Exit(-1)
			}
		}
	})

	flag.Parse()

	mux := server.SetUpServer()

	utils.SetUpLogger()
	utils.AppInfoLogger.Println("Starting up proj in port: " + strconv.Itoa(port))

	portString := strconv.Itoa(port)
	err := http.ListenAndServe(":"+portString, mux)
	if err != nil {
		utils.AppErrorLogger.Println("Failure to bind to the port ", err)
		os.Exit(-2)
	}
}

// -----------------------
//|Set Up server functions|
// -----------------------
func SetUpServer() *http.ServeMux {
	mux := http.NewServeMux()
	data := SetUpData()
	mux.HandleFunc(PingPath, Ping)
	mux.HandleFunc(ShutdownPath, Shutdown)
	mux.Handle(StorePath, data)
	mux.Handle(strings.TrimRight(StorePath, "/"), data)
	mux.Handle(ListPath, data)
	mux.Handle(ListPath+"/", data)
	//mux.Handle("",data)
	return mux
}

func Shutdown(writer http.ResponseWriter, request *http.Request) {
	utils.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
	auth := request.Header.Get("Authorization")
	if auth == Admin {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))

		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(0)
		}()

	} else {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}
}

func Ping(writer http.ResponseWriter, request *http.Request) {
	utils.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("pong"))
}

func (s *storeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	utils.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
	auth := request.Header.Get("Authorization")

	switch {
	case request.Method == http.MethodGet && listBookRe.MatchString(request.URL.Path) || listBooksRe.MatchString(request.URL.Path):
		s.List(writer, request, auth)
		return
	case request.Method == http.MethodGet && StoreRe.MatchString(request.URL.Path):
		s.Get(writer, request, auth)
		return
	case request.Method == http.MethodPut && StoreRe.MatchString(request.URL.Path):
		s.CreateOrUpdate(writer, request, auth)
		return
	case request.Method == http.MethodDelete && StoreRe.MatchString(request.URL.Path):
		s.Delete(writer, request, auth)
		return
	default:
		notFound(writer, request)
		return
	}
}

//List lists all books or given an isbn in path only the one that matches
// Not sure if this method really need auth
func (s *storeHandler) List(writer http.ResponseWriter, request *http.Request, auth string) {
	key := strings.TrimPrefix(request.URL.Path, ListPath)
	key = strings.TrimLeft(key, "/")
	if key == "" {
		//list all books
		s.store.Lock()
		defer s.store.Unlock()
		books := make([]ListInfo, 0, len(s.store.books))
		for key, v := range s.store.books {
			//if v.Owner == auth || auth == Admin{
			books = append(books, ListInfo{
				key,
				v.Owner,
			})
			//}
		}
		jsonBytes, err := json.Marshal(books)
		if err != nil {
			internalServerError(writer, request)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonBytes)
	} else {
		s.store.Lock()
		defer s.store.Unlock()
		book, ok := s.store.books[key]
		if !ok {
			utils.AppErrorLogger.Println("404 key not found")
			msg := "404 key not found"
			http.Error(writer, msg, http.StatusNotFound)
			return
		}
		//if book.Owner == auth || auth == Admin{
		bookList := ListInfo{
			key,
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

//Get given an isbn in path look for it and return full object
func (s *storeHandler) Get(writer http.ResponseWriter, request *http.Request, auth string) {
	key := strings.TrimPrefix(request.URL.Path, StorePath)
	key = strings.TrimLeft(key, "/")
	s.store.Lock()
	book, ok := s.store.books[key]
	s.store.Unlock()
	if !ok {
		utils.AppErrorLogger.Println("Book not found")
		msg := "404 key not found"
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	if book.Owner == auth {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(book.Value))
	} else {
		utils.AppErrorLogger.Println("forbiden", http.StatusForbidden)
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}

}

//Delete given an isbn in path looks for it and if exists deletes it
func (s *storeHandler) Delete(writer http.ResponseWriter, request *http.Request, auth string) {
	key := strings.TrimPrefix(request.URL.Path, StorePath)
	key = strings.TrimLeft(key, "/")
	s.store.Lock()
	book, ok := s.store.books[key]
	s.store.Unlock()
	if !ok {
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
	} else {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
	}

}

//CreateOrUpdate creates entry if isbn doesn't already exist, if exists then updates entry
func (s *storeHandler) CreateOrUpdate(writer http.ResponseWriter, request *http.Request, auth string) {
	var b Book
	key := strings.TrimPrefix(request.URL.Path, StorePath)
	key = strings.TrimLeft(key, "/")
	value, _ := ioutil.ReadAll(request.Body)
	s.store.Lock()
	toUpdate, ok := s.store.books[key]
	s.store.Unlock()
	//ISBN doesn't exist so we create a new entry
	if !ok {
		//Creation part
		s.store.Lock()
		b.Value = string(value)
		b.Owner = auth
		s.store.books[key] = b
		s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	} else {
		//Isbn exists, so we update instead of creating a new one
		if toUpdate.Owner == auth || auth == Admin {
			s.store.Lock()
			updateModel := Book{
				string(value),
				toUpdate.Owner,
			}
			s.store.books[key] = updateModel
			s.store.Unlock()
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("Ok"))
		} else {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Forbidden"))
		}

	}
}
