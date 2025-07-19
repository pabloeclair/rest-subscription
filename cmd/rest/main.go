package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/db"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	shutdownDuration time.Duration
	postgresTimeout  time.Duration
	serverAddr       string
	isSuccessfulInit bool = false
	RedString             = color.New(color.FgRed).SprintFunc()
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

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_USER' does not found"))
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_PASSWORD' does not found"))
	}

	postgresDatabase := os.Getenv("POSTGRES_DATABASE")
	if postgresDatabase == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_DATABASE' does not found"))
	}

	postgresTimeoutString := os.Getenv("SHUTDOWN_DURATION")
	postgresTimeoutInt, err := strconv.Atoi(postgresTimeoutString)
	if err != nil || postgresTimeoutInt <= 0 {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_TIMEOUT' must be a positive number"))
	} else {
		postgresTimeout = time.Duration(postgresTimeoutInt) * time.Second
	}

	if errs != nil {
		color.Red(errs.Error())
		return
	} else {
		isSuccessfulInit = true
		serverAddr = os.Args[1]
		db.DSN = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
			postgresUser, postgresDatabase, postgresPassword, serverAddr)
	}
}

func main() {

	if !isSuccessfulInit {
		return
	}

	fmt.Println("Please, wait...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// connecton to db
	<-time.After(postgresTimeout)
	db, err := gorm.Open(postgres.Open(db.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(RedString("ERROR: connection to postgres: ", err.Error()))
	}

	db.AutoMigrate(&sbscrb.SubscribeDto{})
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()

	// starting server
	mux := http.NewServeMux()
	s := &http.Server{
		Handler: rest.LoggingMiddleware(mux),
		Addr:    serverAddr,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(RedString("ERROR: ", err.Error()))
		}
	}()

	// shutting down server
	log.Printf("Server starting on %s", serverAddr)
	<-ctx.Done()

	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalln(RedString("ERROR: shutdown: ", err.Error()))
	}

}
