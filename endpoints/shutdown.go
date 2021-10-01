package endpoints

import (
	"Week2Proj/constants"
	"Week2Proj/logger"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Shutdown(writer http.ResponseWriter, request *http.Request) {
	logger.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	auth := request.Header.Get("Authorization")
	if auth == constants.Admin {
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
