package endpoints

import (
	"Week2Proj/logger"
	"fmt"
	"net/http"
	"time"
)

func Ping(writer http.ResponseWriter, request *http.Request) {
	logger.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("pong"))
}
