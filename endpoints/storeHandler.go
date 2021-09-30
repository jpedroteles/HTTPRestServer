package endpoints

import (
	"Week2Proj/Types"
	"Week2Proj/logger"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var (
	listBookRe  = regexp.MustCompile(`^\/list[\/]*$`)
	listBooksRe = regexp.MustCompile(`^\/list\/(.+)$`)
	StoreRe     = regexp.MustCompile(`^\/store\/(.+)$`)
)

// StoreHandler handler interface implementation
type StoreHandler struct {
	Store *LocalTypes.KvStore
}

func (s *StoreHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	logger.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
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
		NotFound(writer, request)
		return
	}
}

func (s *StoreHandler) List(writer http.ResponseWriter, request *http.Request, auth string) {
	List(writer, request, auth, s)
}

func (s *StoreHandler) Get(writer http.ResponseWriter, request *http.Request, auth string) {
	Get(writer, request, auth, s)
}

func (s *StoreHandler) CreateOrUpdate(writer http.ResponseWriter, request *http.Request, auth string) {
	CreateOrUpdate(writer, request, auth, s)
}

func (s *StoreHandler) Delete(writer http.ResponseWriter, request *http.Request, auth string) {
	Delete(writer, request, auth, s)
}
