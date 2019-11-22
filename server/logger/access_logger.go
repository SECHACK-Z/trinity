package logger

import (
	"fmt"
	"log"
	"main/pubsub"
	"net/http/httputil"
	"os"
	"time"
)

type LogType struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
}

func fileLogger(event pubsub.Access) {
	req := event.Req
	res := event.Res
	elapsed := event.Elapsed
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		body = []byte("")
	}

	format := "time:%v\tmethod:%v\turi:%v\tstatus:%v\tsize:%v\tapptime:%v\n"
	logData := fmt.Sprintf(format, time.Now().Format("2006-01-02T15:04:05+09:00"), req.Method, req.Host+req.RequestURI, res.StatusCode, len(body), elapsed.Nanoseconds()/1000)
	file, err := os.OpenFile(`./accessLog`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(logData)
}

func StartAccessLogger() {
	pubsub.AccessEvent.Sub(fileLogger)
	// log.Println("AccessLog", req, res, elapsed)
}
