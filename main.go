package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func Serve() {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = ":8081"
	}

	rq := &httputil.ReverseProxy{
		Director: director,
	}

	s := http.Server{
		Addr:    ":8080",
		Handler: rq,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
