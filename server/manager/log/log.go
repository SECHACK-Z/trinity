package log

import (
	"fmt"
	"io"
	"log"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http/httputil"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

type LogManager struct {
	db        *gorm.DB
	systemLog *os.File
	accessLog *os.File
}

func New(db *gorm.DB) *LogManager {
	db.AutoMigrate()
	systemLogFile, err := os.OpenFile("./systemLog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	accessLogFile, err := os.OpenFile(`./accessLog`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	logManager := &LogManager{db: db, systemLog: systemLogFile, accessLog: accessLogFile}
	log.SetOutput(io.MultiWriter(logManager.systemLog, os.Stdout))

	pubsub.AccessEvent.Sub(logManager.AccessLogger)
	pubsub.SystemEvent.Sub(logManager.SystemLogger)

	return logManager
}

func (m *LogManager) AccessLogger(event pubsub.Access) {
	req := event.Req
	res := event.Res
	elapsed := event.Elapsed
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		body = []byte("")
	}

	format := "time:%v\tmethod:%v\turi:%v\tstatus:%v\tsize:%v\tapptime:%v\thost:%v\n"
	logData := fmt.Sprintf(format, time.Now().Format("2006-01-02T15:04:05+09:00"), req.Method, req.Host+req.RequestURI, res.StatusCode, len(body), elapsed.Nanoseconds()/1000, req.Host)
	file, err := os.OpenFile(`./accessLog`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(logData)
}

func (m *LogManager) SystemLogger(event pubsub.System) {
	log.Println(event.Type)
	if event.Type == systemevent.BUILD_FAILED {
		file, err := os.OpenFile(`./buildLog`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		file.WriteString(event.Message)
	}
}
