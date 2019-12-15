package main

import (
	"fmt"
	"net/http"
)

const (
	serverPort = "8080"
	version    = "v0.1.0"
)

type server struct {
	router *http.ServeMux
}

func main() {

	srv := server{router: http.DefaultServeMux}
	srv.routes()

	fmt.Printf("Launching server on port %s\n", serverPort)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		panic(err)
	}
}
