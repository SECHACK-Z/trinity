package main

import (
	"net/http"
	"os"
)

var (
	message = os.Getenv("MESSAGE")
)


func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

// Handler returns port and requested path
func Handler(w http.ResponseWriter, req *http.Request) {
	println("----> :", req.URL.String())
	w.Write([]byte(message + "\n"))
}
