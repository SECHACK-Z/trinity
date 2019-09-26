package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	port string
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	port = os.Args[1]
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":"+port, nil)
}

// Handler returns port and requested path
func Handler(w http.ResponseWriter, req *http.Request) {
	println("----> :"+os.Args[1], req.URL.String())
	w.Write([]byte(port))
}
