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
	_ "github.com/joho/godotenv/autoload"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	shutdownDuration time.Duration
	postgresTimeout  time.Duration
	serverAddr       string
	isSuccessfulInit bool = false
)

func init() {
	var errs error
	if len(os.Args) != 2 {
		errs = errors.New("ERROR: please provide only one argument - server address (e.g. 'localhost:8080')")
	}

	shutdownDurationString := os.Getenv("SHUTDOWN_DURATION")
	shutdownDurationInt, err := strconv.Atoi(shutdownDurationString)
	if err != nil || shutdownDurationInt <= 0 {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'SHUTDOWN_DURATION' must be a positive number"))
	} else {
		shutdownDuration = time.Duration(shutdownDurationInt) * time.Second
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_USER' is not found"))
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_PASSWORD' is not found"))
	}

	postgresDatabase := os.Getenv("POSTGRES_DB")
	if postgresDatabase == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_DB' is not found"))
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_HOST' is not found"))
	}

	postgresTimeoutString := os.Getenv("POSTGRES_TIMEOUT")
	postgresTimeoutInt, err := strconv.Atoi(postgresTimeoutString)
	if err != nil || postgresTimeoutInt <= 0 {
		errs = errors.Join(errs, errors.New("ERROR: the environment variable 'POSTGRES_TIMEOUT' must be a positive number"))
	} else {
		postgresTimeout = time.Duration(postgresTimeoutInt) * time.Second
	}

	models.NotificationInternalError = os.Getenv("NOTIFICATION_INTERNAL_ERROR")
	if models.NotificationInternalError == "" {
		color.Yellow("WARN: the environment variable 'NOTIFICATION_INTERNAL_ERROR' is not found. " +
			"This variable is optional, but you may want to send an extra notification when catching INTERNAL SERVER ERROR " +
			"(e.g., 'Please notify the administrator.')")
	}

	if errs != nil {
		color.Red(errs.Error())
		return
	} else {
		isSuccessfulInit = true
		serverAddr = os.Args[1]
		rest.DSN = fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
			postgresHost, postgresUser, postgresDatabase, postgresPassword)
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
	db, err := gorm.Open(postgres.Open(rest.DSN), &gorm.Config{})
	if err != nil {
		return
	}

	db.AutoMigrate(&models.SubscribeDto{})
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}()

	// starting server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/subscribes", rest.Create)
	mux.HandleFunc("GET /api/v1/subscribes/{id}", rest.Get)
	mux.HandleFunc("GET /api/v1/subscribe", rest.List)
	mux.HandleFunc("PUT /api/v1/subscribes/{id}", rest.Update)
	mux.HandleFunc("PATCH /api/v1/subscribes/{id}", rest.Update)
	mux.HandleFunc("DELETE /api/v1/subscribes/{id}", rest.Delete)
	mux.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		errDto := models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The requested URL %s is not found", r.URL.Path),
			"",
		)
		errDto.Write(w)
	})

	s := &http.Server{
		Handler: rest.LoggingMiddleware(mux),
		Addr:    serverAddr,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(models.RedString("ERROR: ", err.Error()))
		}
	}()

	// shutting down server
	log.Printf("Server starting on %s", serverAddr)
	<-ctx.Done()

	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatal(models.RedString("ERROR: shutdown: ", err.Error()))
	}

}
