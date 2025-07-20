package rest

import (
	"encoding/json"
	"net/http"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DSN string

func Create(w http.ResponseWriter, r *http.Request) {

	var (
		subscribeDto models.SubscribeDto
		errDto       models.FullExceptionDto
	)

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

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{Logger: logger.Default})
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to connect to the database",
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
