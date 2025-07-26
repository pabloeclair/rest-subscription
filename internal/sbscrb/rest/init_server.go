package rest

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
)

func Init() (time.Duration, time.Duration, string, error) {
	var (
		errs             error
		shutdownDuration time.Duration
		postgresTimeout  time.Duration
		serverAddrs      string
	)

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
		return shutdownDuration, postgresTimeout, serverAddrs, errs
	}
	serverAddrs = os.Args[1]
	DSN = fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		postgresHost, postgresUser, postgresDatabase, postgresPassword)
	return shutdownDuration, postgresTimeout, serverAddrs, nil
}
