package endpoints

import (
	"Week2Proj/logger"
	"net/http"
)

func NotFound(writer http.ResponseWriter, request *http.Request) {
	logger.AppErrorLogger.Println("Not found")
	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte(http.StatusText(http.StatusNotFound)))
}
