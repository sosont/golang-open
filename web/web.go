package main

import (
	"fmt"
	"net/http"
)

func helloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, " hello ! this's a http request \n method %v \n request url is %v \n", r.Method, r.URL.String())
}

func panicHandle(w http.ResponseWriter, r *http.Request) {
	panic("version 0.0.1!")
}

func main() {
	// create middleware server
	s := new(MiddlewareServe)
	route := http.NewServeMux()

	route.Handle("/hello", http.HandlerFunc(helloHandle))
	route.Handle("/version", http.HandlerFunc(panicHandle))

	s.Handler = route
	s.Use(LogRequest, ErrCatch)
	// start server
	fmt.Println(http.ListenAndServe(":7000", s))
}