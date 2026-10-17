package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/brightloop/deskbook/internal/api"
	"github.com/brightloop/deskbook/internal/bookings"
	"github.com/brightloop/deskbook/internal/rooms"
)

func main() {
	cfgPath := flag.String("config", "deskbook.json", "room configuration")
	addr := flag.String("addr", "127.0.0.1:8080", "listen address; the SSO proxy forwards to it")
	flag.Parse()

	cfg, err := rooms.Load(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:              *addr,
		Handler:           api.Handler(cfg, bookings.NewStore(), time.Now),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("deskbook listening on %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
