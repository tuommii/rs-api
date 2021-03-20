package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"miikka.xyz/rs/config"
	"miikka.xyz/rs/server"
	"miikka.xyz/rs/store"
)

func main() {
	// Load config from the path provided by the user
	envPath := flag.String("env", ".env", "Path to environment file")
	flag.Parse()
	config := config.Load(*envPath)

	// Init database connection
	database := store.NewDB(config)
	defer database.Close()

	// Init and start server
	s := server.NewServer(database, config)
	go func() {
		s.Start()
	}()
	shutdown(s)
}

func shutdown(s *server.Server) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	log.Println("Exiting...")
	os.Exit(0)
}
