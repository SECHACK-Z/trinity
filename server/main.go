package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	_ "main/statik"

	"github.com/rakyll/statik/fs"
	"github.com/ymotongpoo/goltsv"
)

type Target struct {
	Proxy   string
	Host    string
	Default bool
}

type Config struct {
	Targets []Target
}

func accessLog(req *http.Request) {
	log.Println(req)
	remoteAddr := strings.Split(req.RemoteAddr, ":")[0]
	logData := []map[string]string{{
		"remote_addr":    remoteAddr,
		"request_method": req.Method,
		"request_uri":    req.RequestURI,
		// "https":           "aa",
		// "uri":             "aa",
		// "query_strings":   "aa",
		// "status":          "aa",
		// "bytes_sent":      "aa",
		// "body_bytes_sent": "aa",
	}}

	b := &bytes.Buffer{}
	writer := goltsv.NewWriter(b)
	err := writer.WriteAll(logData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b.String())
}

// NewMultipleHostReverseProxy creates a reverse proxy that will randomly
// select a host from the passed `targets`
func NewMultipleHostReverseProxy(config Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		for _, target := range config.Targets {
			if req.Host == target.Host {
				req.URL.Scheme = "http"
				req.URL.Host = target.Proxy
				//log.Println(req)
				log.Printf("proxy to %s\n", target.Proxy)
				return
			}
		}

		for _, target := range config.Targets {
			if target.Default {
				req.URL.Scheme = "http"
				req.URL.Host = target.Proxy
				accessLog(req)
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
	go func() {
		statikFs, err := fs.New()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		http.Handle("/", http.FileServer(statikFs))
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	log.Fatal(http.ListenAndServe(":9090", proxy))
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler")
		log.Println(r)
		w.Header().Set("X-Ben", "Rad")
		p.ServeHTTP(w, r)
	}
}
