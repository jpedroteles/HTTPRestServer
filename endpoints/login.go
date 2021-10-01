package endpoints

import (
	"Week2Proj/Utils"
	"Week2Proj/logger"
	"fmt"
	"net/http"
	"time"
)

var users = map[string]string{
	"user_a": "passwordA",
	"user_b": "passwordB",
	"user_c": "passwordC",
	"admin":  "Password1",
}

func Login(writer http.ResponseWriter, request *http.Request) {
	logger.HTTPInfoLogger.Println(fmt.Sprintf("Date: %s,IP source: %s,HTTP method: %s,URL: %s", time.Now().Format("2006.01.02 15:04:05"), request.Header.Get("X-FORWARDED-FOR"), request.Method, request.URL))
	writer.Header().Set("content-type", "text/plain; charset=utf-8")
	username, password, _ := request.BasicAuth()

	if users[username] != password {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unauthorized"))
		return
	}
	jwt, err := Utils.CreateJWTPayload(username)
	if err != nil {
		logger.AppErrorLogger.Println("Error Creating jwt token payload")
		internalServerError(writer, request)
		return
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(jwt.Access))
	}
}
