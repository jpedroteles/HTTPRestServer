package endpoints

import "net/http"

func internalServerError(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("internal server error"))
}
