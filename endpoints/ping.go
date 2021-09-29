package endpoints

import (
	"Week2Proj/utils"
	"fmt"
	"net/http"
	"time"
)

func Ping(writer http.ResponseWriter, request *http.Request) {
	utils.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"),request.Header.Get("X-FORWARDED-FOR"),request.Method,request.URL))

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("pong"))
}