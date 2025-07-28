package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	shutdownDuration, postgresTimeout, serverAddrs, err := rest.Init()
	if err != nil {
		color.Red(err.Error())
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
	mux.HandleFunc("GET /api/v1/subscribes/{id}", rest.GetById)
	mux.HandleFunc("GET /api/v1/subscribe", rest.GetList)
	mux.HandleFunc("PUT /api/v1/subscribes/{id}", rest.UpdatePut)
	mux.HandleFunc("PATCH /api/v1/subscribes/{id}", rest.UpdatePatch)
	mux.HandleFunc("DELETE /api/v1/subscribes/{id}", rest.Delete)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		errDto := models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The requested URL %s is not found", r.URL.Path),
			"",
		)
		errDto.Write(w)
	})

	s := &http.Server{
		Handler: rest.LoggingMiddleware(mux),
		Addr:    serverAddrs,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(models.RedString("ERROR: ", err.Error()))
		}
	}()

	// shutting down server
	log.Printf("Server starting on %s", serverAddrs)
	<-ctx.Done()

	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatal(models.RedString("ERROR: shutdown: ", err.Error()))
	}

}
