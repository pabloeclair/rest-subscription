package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN string

func connectToDB(w http.ResponseWriter) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
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

func Get(w http.ResponseWriter, r *http.Request) {

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

func List(w http.ResponseWriter, r *http.Request) {

	var (
		subscribes []models.SubscribeDto
		errDto     models.FullExceptionDto
	)

	db, err := connectToDB(w)
	if err != nil {
		return
	}

	queryParams := r.URL.Query()
	sort := strings.ToUpper(queryParams.Get("sort"))
	value := queryParams.Get("value")

	if (sort == "" && value != "") || (sort != "" && value == "") {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"You should fill in both parameters - 'sort' and 'value'",
			"",
		)
		errDto.Write(w)
		return
	}

	switch sort {
	case "":
		db.Find(&subscribes)
	case "USER_ID":
		db.Where(&models.SubscribeDto{UserId: value}).Find(&subscribes)
	case "SERVICE_NAME":
		db.Where(&models.SubscribeDto{ServiceName: value}).Find(&subscribes)
	default:
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The sort parameter can only be empty or have the values SERVICE_NAME and USER_ID",
			"",
		)
		errDto.Write(w)
		return
	}

	b, err := json.Marshal(&subscribes)
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to marshal response",
			"",
		)
	}

	w.Write(b)
}
