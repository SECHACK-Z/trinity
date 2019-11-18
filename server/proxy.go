package main

import (
	"log"
	"net/http"
	"time"
)

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
