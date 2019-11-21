package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v2"

	_ "main/statik"

	"github.com/labstack/echo/v4"
	"github.com/rakyll/statik/fs"
	"github.com/ymotongpoo/goltsv"
)

type Config struct {
	Targets []Target `json:"targets"`
}

var (
	Logs       []LogType
	resetCh    chan struct{}
	httpsHosts []string
	config     Config
)

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
	systemLogFile, err := os.OpenFile("./systemLog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("cannot open systemlog")
	}
	defer systemLogFile.Close()
	log.SetOutput(io.MultiWriter(systemLogFile, os.Stdout))

	log.Println("Server Start.")
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
	log.Printf("%d directors registrated.\n", len(config.Targets))

	// Web UI
	go func() {
		e := echo.New()
		e.HideBanner = true
		log.Println("webUI started.")
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

		e.GET("/api/config", func(c echo.Context) error {
			buf, err := yaml.Marshal(config)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return c.JSON(http.StatusOK, struct {
				Yaml string
			}{
				Yaml: string(buf),
			})
		})

		e.POST("/api/config", func(c echo.Context) error {
			req := &struct {
				Yaml string `json:"yaml"`
			}{}
			c.Bind(req)

			newConfig := Config{}
			if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}

			config = newConfig
			resetCh <- struct{}{}

			return c.NoContent(http.StatusOK)
		})

		e.POST("/api/config/save", func(c echo.Context) error {
			req := &struct {
				Name string `json:"name"`
				Yaml string `json:"yaml"`
			}{}
			c.Bind(req)

			newConfig := Config{}
			if err := yaml.Unmarshal([]byte(req.Yaml), &newConfig); err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}

			config = newConfig

			if req.Name == "" {
				req.Name = "config.yaml"
			}

			if err := ioutil.WriteFile(req.Name, []byte(req.Yaml), 0755); err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			log.Println("New Configuration file saved.")
			resetCh <- struct{}{}

			return c.NoContent(http.StatusOK)
		})

		if err := e.Start(":8080"); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// 設定再読み込みなどなどを行う
	for {
		resetCh = make(chan struct{})
		httpsHosts = make([]string, 0)

		proxy := NewMultipleHostReverseProxy(config)
		proxy.Transport = &myTransport{}

		httpsSrv := &http.Server{Handler: proxy}
		httpSrv := &http.Server{Addr: ":80", Handler: proxy}

		for _, target := range config.Targets {
			if target.Https {
				httpsHosts = append(httpsHosts, target.Host)
			}
		}
		log.Println("New settings applied.")
		go func() {
			httpsSrv.Serve(autocert.NewListener(httpsHosts...))
		}()
		go func() {
			httpSrv.ListenAndServe()
		}()
		<-resetCh

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpsSrv.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()

		go func() {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
		wg.Wait()
	}

}
