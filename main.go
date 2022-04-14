package main

import (
	"log"
	"net/http"
	"time"

	"github.com/os/ratelimiter/rate"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		log.Printf("Failed to return a response: %s", err)
	}
}

func main() {
	store := rate.NewMemoryStore(time.Second * 5)
	limiter := rate.NewFixedWindowLimiter(15, time.Second*6, store)
	burstLimiter := rate.NewFixedWindowLimiter(10, time.Second*3, store)
	identifier := rate.NewIPIdentifier()
	limit := rate.Limit(identifier, limiter, burstLimiter)
	http.HandleFunc("/hello", limit(handleHello))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
