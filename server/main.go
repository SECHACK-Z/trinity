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

var Logs []LogType

func accessLog(res *http.Request) error {

	//format := "{method:\"%v\", uri:\"%v\"},\n"
	l := LogType{
		Method: res.Method,
		URI:    res.RequestURI,
	}
	Logs = append(Logs, l)

	format := "time:%v\tmethod:%v\turi:%v\tstatus:200\tsize:10\tapptime:0.100\n"
	logData := fmt.Sprintf(format, time.Now().Format("2006-01-02T15:04:05+09:00"), res.Method, res.RequestURI)
	fmt.Println(logData)
	file, err := os.OpenFile(`./logFile`, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(logData)
	return nil
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
		e := echo.New()
		statikFs, err := fs.New()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		e.GET("/api/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		e.GET("/", echo.WrapHandler(http.FileServer(statikFs)))
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
