package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	"github.com/stdatiks/go-daemon"
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

func configure() *daemon.Config {
	s := &service{}
	s.server = &http.Server{
		Addr:           ":8080",
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	cfg := daemon.NewConfig()
	cfg.Name = "cmd_example"
	cfg.Description = "Command Line daemon example"
	cfg.StartHdlr = s.Start
	cfg.StopHdlr = s.Stop
	cfg.HideMethodsWarning = true
	return cfg
}

func main() {
	d, _ := daemon.New(configure())

	var (
		usage  = "Usage: " + os.Args[0] + " install | uninstall | restart | start | stop | status"
		status string
		err    error
	)
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			status, err = d.Install(os.Args[2:]...)
		case "uninstall":
			status, err = d.Uninstall()
		case "restart":
			status, err = d.Restart()
		case "start":
			status, err = d.Start()
		case "stop":
			status, err = d.Stop()
		case "status":
			status, err = d.Status()
		default:
			err = fmt.Errorf(usage)
		}
	} else {
		err = d.Run()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if status != "" {
		fmt.Println(status)
	}
}
