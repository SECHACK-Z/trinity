package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	_ "main/statik"

	"github.com/rakyll/statik/fs"
)

type Target struct {
	Proxy   string
	Host    string
	Default bool
}

type Config struct {
	Targets []Target
}

func accessLog(res *http.Request) error {
	log.Println(res)
	//remoteAddr := strings.Split(req.RemoteAddr, ":")[0]

	// time:2015-09-06T05:58:05+09:00	method:POST	uri:/foo/bar?token=xxx&uuid=1234	status:200	size:12	apptime:0.057
	// logData := []map[string]string{{
	// 	"time": time.Now().Format("2006-01-02T15:04:05+09:00"),
	// 	//"remote_addr": remoteAddr,
	// 	"method":      req.Method,
	// 	"request_uri": req.RequestURI,
	// 	"status":      "200",
	// 	"size":        string(req.ContentLength),
	// 	"apptime":     "-",
	// 	// "uri":             "aa",
	// 	// "query_strings":   "aa",
	// 	// "status":          "aa",
	// 	// "bytes_sent":      "aa",
	// 	// "body_bytes_sent": "aa",
	// }}
	format := "time:%v\tmethod:%v\turi:%v\tstatus:200\tsize:%v\tapptime:0.100\n"
	logData := fmt.Sprintf(format, time.Now().Format("2006-01-02T15:04:05+09:00"), res.Method, res.RequestURI, len(res.RequestURI))

	// b := &bytes.Buffer{}
	// writer := goltsv.NewWriter(b)
	// err := writer.WriteAll(logData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(b.String())
	file, err := os.OpenFile(`./logFile`, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(logData)
	return nil
}

// func modifier(res *http.Response) error {
// 	fmt.Println(res)
// 	return nil
// }

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
	// return &httputil.ReverseProxy{Director: director, ModifyResponse: modifier}
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
