package main

import (
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	_ "main/statik"

	"github.com/rakyll/statik/fs"
	"github.com/labstack/echo/v4"
)

type Target struct {
	Proxy string
	Host string
	Https bool
	ForceHttps bool
	Default bool
}

type Config struct {
	Targets []Target
}

// NewMultipleHostReverseProxy creates a reverse proxy that will randomly
// select a host from the passed `targets`
func NewMultipleHostReverseProxy(config Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		for _, target := range config.Targets {
			if req.Host == target.Host {
				req.URL.Scheme = "http"
				req.URL.Host = target.Proxy
				log.Printf("proxy to %s\n", target.Proxy)
				return
			}
		}


		for _, target := range config.Targets {
			if target.Default {
				req.URL.Scheme = "http"
				req.URL.Host = target.Proxy
				log.Printf("proxy to %s\n", target.Proxy)
				return
			}
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	var config Config
	body, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("file read Error: %v", err)
	}

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		log.Fatalf("yaml parse Error: %v", err)
	}

	proxy := NewMultipleHostReverseProxy(config)
	go func () {
		e := echo.New()
		statikFs, err := fs.New()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		e.GET("/api/ping", func (c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		e.GET("/", echo.WrapHandler(http.FileServer(statikFs)))

		if err := e.Start(":8080"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	httpsHosts := make([]string, 0)
	for _, target := range config.Targets {
		if target.Https {
			httpsHosts = append(httpsHosts, target.Host)
		}
	}

	go func () {
		log.Fatal(http.Serve(autocert.NewListener(httpsHosts...), proxy))
	}()
	log.Fatal(http.ListenAndServe(":80", proxy))
}
