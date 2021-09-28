package endpoints

import (
	"net/http"
)

func Ping(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","application/json")
	switch  request.Method{
	case http.MethodGet:
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Pong"))
		return
	default:
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"message": "not found"}`))
	}
}