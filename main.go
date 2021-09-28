package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"utils"
)
const (

	//Root     = "/api/"
	StorePath = "/store/"
	ListPath  ="/list"
	pingPath  = "/ping"
)

var(
	listBookRe   = regexp.MustCompile(`^\/list[\/]*$`)
	listBooksRe    = regexp.MustCompile(`^\/list\/(\d+)$`)
	StoreRe    = regexp.MustCompile(`^\/store\/(\d+)$`)
	CreateBookRe   = regexp.MustCompile(`^\/store[\/]*$`)
)


// ------
// |TYPES|
// ------

//Book type definitions
type Book struct {
	ISBN      int    `json:"isbn"`
	Title  	  string `json:"title"`
	Author 	  string `json:"author"`
	Owner string `json:"owner"`
}

//kvStore type definition. Mutex for lock/unlock when making operations on object
type kvStore struct {
	books map[string]Book
	*sync.RWMutex
}

//storeHandler handler interface implementation
type storeHandler struct {
	store *kvStore
}

type listInfo struct{
	key string `json:"key"`
	owner string `json:"owner"`
}

// NewBook constructor
func NewBook(isbn int, title string, author string) *Book{
	Book:= Book{ISBN:isbn, Title: title, Author: author}
	return &Book
}

// Stringer interface
func (b *Book) String() string {
	return fmt.Sprintf("Book(%d, %s by %s)", b.ISBN, b.Title,
		b.Author)
}

func SetUpData() *storeHandler {
	book := &storeHandler{
		store : &kvStore{
			books:map[string]Book{
				"1":{1,"Get Set Go!", "John Smith", "Test"},
				"2":{2,"Be a Go Getter", "John Smith", "David Byrne"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	return book
}


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
				utils.AppErrorLogger.Println("Failure to parte to int ", err)
				os.Exit(-1)
			}
		}
	})

	flag.Parse()

	utils.SetUpLogger()
	utils.AppInfoLogger.Println("Starting up proj in port: "+ strconv.Itoa(port))

	mux := SetUpServer()

	//fmt.Println("Server Available - see")
	//fmt.Println("\t", fmt.Sprintf("http://%s:%s%s", ConnHost, port, ping))

	portString:= strconv.Itoa(port)
	err:= http.ListenAndServe(":"+portString, mux)
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
	mux.HandleFunc(pingPath,Ping)
	mux.Handle(StorePath,data)
	mux.Handle(strings.TrimRight(StorePath,"/"),data)
	mux.Handle(ListPath,data)
	mux.Handle(ListPath +"/", data)
	//mux.Handle("",data)
	return mux
}

func Ping(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s storeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	switch {
	case request.Method == http.MethodGet && listBookRe.MatchString(request.URL.Path) || listBooksRe.MatchString(request.URL.Path):
		s.List(writer, request)
		return
	case request.Method == http.MethodGet && StoreRe.MatchString(request.URL.Path):
		s.Get(writer, request)
		return
	case request.Method == http.MethodPut && CreateBookRe.MatchString(request.URL.Path):
		s.CreateOrUpdate(writer, request)
		return
	case request.Method == http.MethodDelete && StoreRe.MatchString(request.URL.Path):
		s.Delete(writer, request)
		return
	default:
		notFound(writer, request)
		return
	}
}

// List TODO check problem with converting object to json
//List list all books or given a isbn in path only the one that matchs
func (s storeHandler) List(writer http.ResponseWriter, request *http.Request) {
	key := strings.TrimPrefix(request.URL.Path,ListPath)
	key = strings.TrimLeft(key,"/")
	if key ==""{
		//list all books
		s.store.RLock()
		books := make([]listInfo, 0, len(s.store.books))
		utils.AppWarningLogger.Println(s.store.books)
		for _, v := range s.store.books {
			books = append(books, listInfo{
				strconv.Itoa(v.ISBN),
				v.Owner,
			})
		}
		s.store.RUnlock()
		jsonBytes, err := json.Marshal(books)
		if err != nil {
			internalServerError(writer, request)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write(jsonBytes)
	}else{
		s.store.RLock()
		book, ok := s.store.books[key]
		s.store.RUnlock()
		if !ok{
			utils.AppErrorLogger.Println("404 key not found")
			msg := "404 key not found"
			http.Error(writer, msg, http.StatusNotFound)
			return
		}
		bookList := listInfo{
			strconv.Itoa(book.ISBN),
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
	}
}

//Get given an isbn in path look for it and return full object
func (s storeHandler) Get(writer http.ResponseWriter, request *http.Request) {
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
	jsonData, err := json.Marshal(book)
	if err != nil {
		internalServerError(writer, request)
		return
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(jsonData))
	}
}

//Delete given an isbn in path looks for it and if exists deletes it
func (s storeHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	key := strings.TrimPrefix(request.URL.Path,StorePath)
	key = strings.TrimLeft(key,"/")
	s.store.RLock()
	_, ok := s.store.books[key]
	s.store.RUnlock()
	if !ok{
		utils.AppErrorLogger.Println("Book not found")
		msg := "404 key not found"
		http.Error(writer, msg, http.StatusNotFound)
		return
	}
	s.store.Lock()
	//Remove key
	delete(s.store.books, key)
	s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
}

//CreateOrUpdate creates entry if isbn doesn't already exist, if exists then updates entry
func (s storeHandler) CreateOrUpdate(writer http.ResponseWriter, request *http.Request) {
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
		s.store.books[strconv.Itoa(b.ISBN)] = b
		s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	}else{
		//Isbn exists, so we update instead of creating a new one
		s.store.Lock()
		updateModel := Book{
			b.ISBN,
			b.Title,
			b.Author,
			b.Owner,
		}
		s.store.books[string(b.ISBN)] = updateModel
		s.store.Unlock()
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Ok"))
	}


}

func internalServerError(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("internal server error"))
}

func notFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("not found"))
}

