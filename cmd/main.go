package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"psui"
	"psui/server"
	"psui/service"
	"time"
)

func main() {

	var err error
	cfg := &psui.Config{}
	err = cfg.Load("config.json")
	if err != nil {
		log.Fatalf("err loading config file: %s", err)
	}

	svc := service.NewService(*cfg)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := server.NewRouter(server.NewHandlers(svc, *cfg))

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	srv.Shutdown(context.Background())
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
