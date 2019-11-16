package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v2"

	_ "main/statik"

	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
	"github.com/ymotongpoo/goltsv"
)

type Target struct {
	Proxy      string
	Host       string
	Https      bool
	ForceHttps bool
	Default    bool
}

type Config struct {
	Targets []Target
}

type LogType struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
}
type myTransport struct{}

var Logs []LogType

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	response, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("Server is not reachable, err")
		return nil, err
	}
	elapsed := time.Since(start)
	log.Println("Response Time:", elapsed.Nanoseconds())

	accessLog(req, response, elapsed)

	return response, nil
}

func accessLog(req *http.Request, res *http.Response, elapsed time.Duration) {
	// log.Println("AccessLog", req, res, elapsed)
	l := LogType{
		Method: req.Method,
		URI:    req.RequestURI,
	}
	Logs = append(Logs, l)
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		body = []byte("")
	}

	format := "time:%v\tmethod:%v\turi:%v\tstatus:%v\tsize:%v\tapptime:%v\n"
	logData := fmt.Sprintf(format, time.Now().Format("2006-01-02T15:04:05+09:00"), req.Method, req.Host+req.RequestURI, res.StatusCode, len(body), elapsed.Milliseconds())
	log.Println(logData)
	file, err := os.OpenFile(`./logFile`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(logData)
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
	proxy.Transport = &myTransport{}
	go func() {
		e := echo.New()
		statikFs, err := fs.New()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		e.GET("/api/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		h := http.FileServer(statikFs)

		e.GET("/*", echo.WrapHandler(http.StripPrefix("/", h)))
		e.GET("/api/log", func(c echo.Context) error {
			//デモ用の実装
			data, err := ioutil.ReadFile(`logFile`)
			if err != nil {
				log.Println(err)
			}
			b := bytes.NewBufferString(string(data))
			reader := goltsv.NewReader(b)
			records, _ := reader.ReadAll()
			bytes, _ := json.Marshal(records)
			return c.JSON(200, string(bytes))
		})
		e.GET("/api/alp", func(c echo.Context) error {
			out, err := exec.Command("sh", "-c", "cat logFile | alp --sort=max ltsv").Output()
			if err != nil {
				log.Fatal(err.Error())
			}
			return c.String(200, string(out))
		})
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

	go func() {
		log.Fatal(http.Serve(autocert.NewListener(httpsHosts...), proxy))
	}()
	log.Fatal(http.ListenAndServe(":80", proxy))
}
