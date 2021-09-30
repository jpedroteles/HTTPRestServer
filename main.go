package main

import (
	"Week2Proj/constants"
	"Week2Proj/logger"
	"Week2Proj/server"
	"flag"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var (
		port int
		depth int
	)

	//Handling arguments
	flag.IntVar(&port, "port", 8000, "port to listen on")
	flag.IntVar(&depth, "depth", 0, "LRU depth (zero for unbounded)")
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			_, err := strconv.Atoi(f.Value.String())
			if err != nil {
				logger.AppErrorLogger.Println("Failure to parse to int ", err)
				os.Exit(-1)
			}
		}
		if f.Name == "depth" {
			_, err := strconv.Atoi(f.Value.String())
			if err != nil {
				logger.AppErrorLogger.Println("Invalid depth value, default depth used(100)", err)
				depth = constants.DefaultDepth
			}
		}
	})

	flag.Parse()

	mux := server.SetUpServer(depth)

	logger.SetUpLogger()
	logger.AppInfoLogger.Println("Starting up proj in port: " + strconv.Itoa(port))

	portString := strconv.Itoa(port)
	err := http.ListenAndServe(":"+portString, mux)
	if err != nil {
		logger.AppErrorLogger.Println("Failure to bind to the port ", err)
		os.Exit(-2)
	}
}
