package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var (
	message = os.Getenv("MESSAGE")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":9091", nil)
}

// Handler returns port and requested path
func Handler(w http.ResponseWriter, req *http.Request) {
	sleepTime := rand.Intn(5)
	println("----> :", req.URL.String())
	println("zzz... for ", sleepTime)
	time.Sleep(time.Second * time.Duration(sleepTime))
	w.Write([]byte(message + "\n"))
	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("%s", dump)
}
