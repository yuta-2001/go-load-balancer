package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"sync"	
	"io/ioutil"
	"encoding/json"
	"net/url"
)

type Config struct {
	Proxy Proxy `json:"proxy"`
	Backends []Backend `json:"backends"`
}

type Proxy struct {
	Port string `json:"port"`
}

type Backend struct {
	URL string `json:"url"`
	IsDead bool
	mu sync.RWMutex
}

var mu sync.Mutex
var idx int = 0

func IbHandler(w http.ResponseWriter, r *http.Request) {
	maxLen := len(cfg.Backends)

	mu.Lock()
	// currentBackend := cfg.Backends[idx%maxLen]
	targetUrl, err := url.Parse(cfg.Backends[idx%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	idx++
	mu.Unlock()
	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)
	reverseProxy.ServeHTTP(w, r)
}

var cfg Config

func Serve() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal(data, &cfg)

	s := http.Server{
		Addr: ":" + cfg.Proxy.Port,
		Handler: http.HandlerFunc(IbHandler),
	}

	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
