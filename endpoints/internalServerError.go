package endpoints

import (
	"Week2Proj/logger"
	"net/http"
)

func internalServerError(writer http.ResponseWriter, request *http.Request) {
	logger.AppErrorLogger.Println("Internal error")
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}
