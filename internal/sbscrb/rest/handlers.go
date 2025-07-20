package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DSN string

func connectToDB(w http.ResponseWriter) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{Logger: logger.Default})
	if err != nil {
		errDto := models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to connect to the database",
			err.Error(),
		)
		errDto.Write(w)
		return nil, err
	}
	return db, nil
}

func Create(w http.ResponseWriter, r *http.Request) {

	var (
		subscribeDto models.SubscribeDto
		errDto       models.FullExceptionDto
	)

	db, err := connectToDB(w)
	if err != nil {
		return
	}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&subscribeDto); err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"Incorrect JSON body",
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	res := db.Create(&subscribeDto)
	if res.Error != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to create subscribe",
			res.Error.Error(),
		)
		errDto.Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetById(w http.ResponseWriter, r *http.Request) {

	var (
		subscribeDto models.SubscribeDto
		errDto       models.FullExceptionDto
	)

	idStr := r.PathValue("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"Incorrect id in the URL path. Please specify a positive number",
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	db, err := connectToDB(w)
	if err != nil {
		return
	}

	res := db.First(&subscribeDto, idInt)
	if res.Error != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	}

	b, err := json.Marshal(&subscribeDto)
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to marshal a response",
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	w.Write(b)
}
