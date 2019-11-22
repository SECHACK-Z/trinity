package transport

import (
	"log"
	"main/pubsub"
	"net/http"
	"time"
)

type myTransport struct{
}

func New() *myTransport {
	return &myTransport{
	}
}

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	response, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("Server is not reachable, err")
		return nil, err
	}
	elapsed := time.Since(start)
	log.Println("Response Time:", elapsed.Nanoseconds())

	pubsub.AccessEvent.Pub(pubsub.Access{req, response, elapsed})
	return response, nil
}
