package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type APIConfig struct {
	Endpoint string
	Scheme string
	Host string
}

type Config struct {
	API []url.URL
}

// NewMultipleHostReverseProxy creates a reverse proxy that will randomly
// select a host from the passed `targets`
func NewMultipleHostReverseProxy(targets []url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target := targets[rand.Int()%len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	var config Config
	
	_, err := toml.DecodeFile("config.toml",&config)
	if err != nil {
		fmt.Errorf("toml decode Error: %v", err)
	}

	proxy := NewMultipleHostReverseProxy(config.API)
	log.Fatal(http.ListenAndServe(":9090", proxy))
}
