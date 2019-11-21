package transport

import (
	"log"
	"main/pubsub"
	"net/http"
	"time"
)

type myTransport struct{
	ps *pubsub.AccessEventPubSub
}

func New() *myTransport {
	return &myTransport{
		ps: pubsub.GetAccessEventPubSub(),
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

	t.ps.Pub(pubsub.AccessEvent{req, response, elapsed})
	return response, nil
}
