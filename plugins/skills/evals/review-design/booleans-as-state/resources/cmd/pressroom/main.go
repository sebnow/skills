package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/brightloop/pressroom/internal/articles"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address; the SSO proxy forwards to it")
	flag.Parse()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           articles.Handler(articles.NewService(time.Now)),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("pressroom listening on %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
