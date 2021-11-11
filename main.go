package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"jobsity-code-challenge/broker"
	"jobsity-code-challenge/config"
	"jobsity-code-challenge/repo"
	"jobsity-code-challenge/server"
	"jobsity-code-challenge/token"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	configs := config.Load()

	rabbit := broker.New(configs.RabbitMQ)
	db := repo.SetupDB(configs.Repo)
	defer rabbit.Close()
	defer db.Close()

	tokenService := token.New(configs.JwtSecretKey)

	router := server.Setup(rabbit, db, tokenService, configs.StooqURLString)

	crt, _ := tls.LoadX509KeyPair(configs.CrtFile, configs.Keyfile)
	srv := &http.Server{
		Addr:    ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{crt},
		},
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Fatal error starting server: %v \n", err))
		}
	}()
	go func() {
		if err := http.ListenAndServe(":80", router); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Fatal error starting server: %v \n", err))
		}
	}()
	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Fatal error shutdown server: %v \n", err))
	}
	log.Println("Finished Server")
}
