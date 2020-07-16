package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/nsemikov/go-daemon"
)

type service struct {
	server *http.Server
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (s service) Start() error {
	// Non-blocking function
	go s.server.ListenAndServe()
	fmt.Println("started at http://localhost:8080")
	return nil
}

func (s service) Stop() error {
	// Non-blocking function
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func main() {
	s := &service{}
	s.server = &http.Server{
		Addr:           ":8080",
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	cfg := daemon.NewConfig()
	cfg.StartHdlr = s.Start
	cfg.StopHdlr = s.Stop
	cfg.HideMethodsWarning = true
	d, _ := daemon.New(cfg)
	d.Run()
}
