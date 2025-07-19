package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
)

var (
	shutdownDuration time.Duration
	serverAddr       string
	isSuccessfulInit bool = false
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	RedString        = color.New(color.FgRed).SprintFunc()
)

func init() {
	var errs error
	if len(os.Args) != 2 {
		errs = errors.New("ERROR: please provide only one argument - server address (e.g. 'localhost:8080')")
	}

	shutdownDurationString := os.Getenv("SHUTDOWN_DURATION")
	shutdownDurationInt, err := strconv.Atoi(shutdownDurationString)
	if err != nil || shutdownDurationInt <= 0 {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'SERVER_SHUTDOWN' must be a positive number"))
	} else {
		shutdownDuration = time.Duration(shutdownDurationInt) * time.Second
	}

	PostgresUser := os.Getenv("POSTGRES_USER")
	if PostgresUser == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_USER' does not found"))
	}

	PostgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if PostgresPassword == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_PASSWORD' does not found"))
	}

	PostgresDatabase := os.Getenv("POSTGRES_DATABASE")
	if PostgresDatabase == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_DATABASE' does not found"))
	}

	if errs != nil {
		color.Red(errs.Error())
		return
	} else {
		isSuccessfulInit = true
		serverAddr = os.Args[1]
	}
}

func main() {

	if !isSuccessfulInit {
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	mux := http.NewServeMux()

	s := &http.Server{
		Handler: mux,
		Addr:    serverAddr,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(RedString("ERROR: ", err))
		}
	}()

	log.Printf("Server starting on %s", serverAddr)
	<-ctx.Done()

	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalln(RedString("ERROR: shutdown: ", err))
	}

}
