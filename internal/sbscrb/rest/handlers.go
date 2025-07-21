package rest

import (
	"encoding/json"
	"errors"
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
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	} else if res.Error != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to get the subscribe with id = %d", idInt),
			res.Error.Error(),
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
			"The 'sort' or 'value' parameters are missing. Please fill in both parameters",
			"",
		)
		errDto.Write(w)
		return
	}

	var res *gorm.DB
	switch sort {
	case "":
		res = db.Find(&subscribes)
	case "USER_ID":
		res = db.Where(&models.SubscribeDto{UserId: value}).Find(&subscribes)
	case "SERVICE_NAME":
		res = db.Where(&models.SubscribeDto{ServiceName: value}).Find(&subscribes)
	default:
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"Incorrect the 'sort' parameter. The 'sort' can only be empty or have the values 'SERVICE_NAME' and 'USER_ID'",
			"",
		)
		errDto.Write(w)
		return
	}

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to get the subscribe list",
			res.Error.Error(),
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

func Delete(w http.ResponseWriter, r *http.Request) {

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

	res := db.Delete(&subscribeDto, idInt)
	if res.Error != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to delete the subscribe with id = %d", idInt),
			res.Error.Error(),
		)
		errDto.Write(w)
		return
	} else if res.RowsAffected == 0 {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
