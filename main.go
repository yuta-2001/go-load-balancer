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

func (backend *Backend) SetDead(b bool) {
	backend.mu.Lock()
	backend.IsDead = b
	backend.mu.Unlock()
}

func (backend *Backend) GetIsDead() bool {
	backend.mu.RLock()
	isAlive := backend.IsDead
	backend.mu.RUnlock()
	return isAlive
}

var mu sync.Mutex
var idx int = 0

func IbHandler(w http.ResponseWriter, r *http.Request) {
	maxLen := len(cfg.Backends)

	mu.Lock()
	currentBackend := cfg.Backends[idx%maxLen]
	if currentBackend.GetIsDead() {
		idx++
	}
	targetUrl, err := url.Parse(cfg.Backends[idx%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	idx++
	mu.Unlock()
	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("%v is dead.", targetUrl)
		currentBackend.SetDead(true)
		IbHandler(w, r)
	}
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
