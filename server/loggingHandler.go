package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type Target struct {
	Proxy      string `json:"proxy"`
	Host       string `json:"host"`
	Https      bool   `json:"https"`
	ForceHttps bool   `json:"forceHttps"`
	Default    bool   `json:"default"`
}

type LogType struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
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
