package server

import (
	"Week2Proj/constants"
	_ "Week2Proj/constants"
	"Week2Proj/endpoints"
	_ "Week2Proj/endpoints"
	"Week2Proj/utils"
	_ "Week2Proj/utils"
	"net/http"
	"strings"
)

func SetUpServer() *http.ServeMux{
	mux := http.NewServeMux()
	data := utils.SetUpData()
	mux.HandleFunc(constants.PingPath,endpoints.Ping)
	mux.HandleFunc(constants.ShutdownPath, endpoints.Shutdown)
	mux.Handle(constants.StorePath,data)
	mux.Handle(strings.TrimRight(constants.StorePath,"/"),data)
	mux.Handle(constants.ListPath,data)
	mux.Handle(constants.ListPath+"/", data)
	return mux
}
