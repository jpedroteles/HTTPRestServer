package server

import "endpoints"

func SetUpServer(){
	//mux := http.NewServeMux()

	//mux.Handle("store", endpoints.Ping())

	endpoints.Ping1()
}
