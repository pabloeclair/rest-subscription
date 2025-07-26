package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pabloeclair/rest-subscription/internal/sbscrb/models"
	"github.com/pabloeclair/rest-subscription/internal/sbscrb/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN string

func connectToDB(w http.ResponseWriter) (*repositories.GormSubscribeRepository, error) {
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
	repo := repositories.GormSubscribeRepository{Db: db}
	return &repo, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	var (
		//the subscribe from request
		subscribeDto models.SubscribeDto
		//the error for response
		errDto models.FullExceptionDto
	)

	// headers validation
	if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The request body must be in JSON format",
			"",
		)
		errDto.Write(w)
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

	// body validation
	if err := subscribeDto.Validate(); err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			err.Error(),
			"",
		)
		errDto.Write(w)
		return
	}

	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	// create operation
	err = repo.Create(subscribeDto.ToDatabase())
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to create the subscribe",
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	// result
	w.WriteHeader(http.StatusCreated)
	w.Header().Del("Content-Type")
}

func GetById(w http.ResponseWriter, r *http.Request) {
	var (
		//the subscribe from db
		subscribeDb *models.Subscribe
		//the error for response
		errDto models.FullExceptionDto
	)

	// query validate
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

	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	// find operation
	subscribeDb, err = repo.FindByID(uint(idInt))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errDto = models.NewFullExceptionDto(
				http.StatusNotFound,
				fmt.Sprintf("The subscribe with id = %d is not found", idInt),
				err.Error(),
			)
			errDto.Write(w)
			return
		} else {
			errDto = models.NewFullExceptionDto(
				http.StatusInternalServerError,
				fmt.Sprintf("Failed to get the subscribe with id = %d", idInt),
				err.Error(),
			)
			errDto.Write(w)
			return
		}
	}

	// result
	b, err := json.Marshal(subscribeDb.ToDto())
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

func GetList(w http.ResponseWriter, r *http.Request) {
	var (
		//the subscribes from db
		subscribes []*models.Subscribe
		//the subscribes for response
		subscribesDto []*models.SubscribeDto
		//the error for response
		errDto models.FullExceptionDto
	)

	queryParams := r.URL.Query()
	sort := strings.ToUpper(queryParams.Get("sort"))
	value := queryParams.Get("value")

	// query validate
	if (sort == "" && value != "") || (sort != "" && value == "") {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The 'sort' or 'value' parameters are missing. Please fill in both parameters",
			"",
		)
		errDto.Write(w)
		return
	}

	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	// find operation
	switch sort {
	case "":
		subscribes, err = repo.FindAll()
	case "USER_ID":
		subscribes, err = repo.FindByUserId(value)
	case "SERVICE_NAME":
		subscribes, err = repo.FindByServiceName(value)
	default:
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"Incorrect the 'sort' parameter. The 'sort' can only be empty or have the values 'SERVICE_NAME' and 'USER_ID'",
			"",
		)
		errDto.Write(w)
		return
	}

	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to get the subscribe list",
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	// result
	for _, v := range subscribes {
		subscribesDto = append(subscribesDto, v.ToDto())
	}

	b, err := json.Marshal(&subscribesDto)
	if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			"Failed to marshal response",
			"",
		)
	}

	w.Write(b)
}

func UpdatePatch(w http.ResponseWriter, r *http.Request) {

	var (
		//the subscribe from db
		subscribeDb *models.Subscribe
		//the subscribe from request
		subscribeDto *models.SubscribeDto
		//the error for response
		errDto models.FullExceptionDto
	)

	// headers validate
	if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The request body must be in JSON format",
			"",
		)
		errDto.Write(w)
		return
	}

	// query validate
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

	// body unmarshal
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

	// fields validate
	if err = subscribeDto.ValidateTime(); err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			err.Error(),
			"",
		)
		errDto.Write(w)
		return
	}

	// preparing fields
	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	subscribeDb, err = repo.FindByID(uint(idInt))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	} else if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to get the subscribe with id = %d", idInt),
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	if subscribeDto.ServiceName != "" {
		subscribeDb.ServiceName = subscribeDto.ServiceName
	}
	if subscribeDto.Price != nil {
		subscribeDb.Price = *subscribeDto.Price
	}
	if subscribeDto.UserId != "" {
		subscribeDb.UserId = subscribeDto.UserId
	}
	if !subscribeDto.StartDate.IsZero() {
		subscribeDb.StartDate = subscribeDto.StartDate
	}
	if subscribeDto.EndDate != nil {
		subscribeDb.EndDate = subscribeDto.EndDate
	}
	if subscribeDto.EndDate != nil && subscribeDto.StartDate.After(*subscribeDto.EndDate) {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The field 'end_time' should be after the 'start_time'",
			"",
		)
		errDto.Write(w)
		return
	}

	// update operation
	err = repo.Update(uint(idInt), subscribeDb)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	} else if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to update the subscribe with id = %d", idInt),
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	// result
	w.WriteHeader(http.StatusNoContent)
	w.Header().Del("Content-Type")

}

func UpdatePut(w http.ResponseWriter, r *http.Request) {
	var (
		//the subscribe from request
		subscribeDto *models.SubscribeDto
		//the error for response
		errDto models.FullExceptionDto
	)

	// headers validate
	if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			"The request body must be in JSON format",
			"",
		)
		errDto.Write(w)
		return
	}

	// query validate
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

	// body unmarshal
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

	// fields validate
	if err = subscribeDto.Validate(); err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusBadRequest,
			err.Error(),
			"",
		)
		errDto.Write(w)
		return
	}

	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	// update operation
	err = repo.Update(uint(idInt), subscribeDto.ToDatabase())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	} else if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to update the subscribe with id = %d", idInt),
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	// result
	w.WriteHeader(http.StatusNoContent)
	w.Header().Del("Content-Type")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var (
		//the error for response
		errDto models.FullExceptionDto
	)

	// query validate
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

	repo, err := connectToDB(w)
	if err != nil {
		return
	}

	// delete operation
	err = repo.Delete(uint(idInt))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errDto = models.NewFullExceptionDto(
			http.StatusNotFound,
			fmt.Sprintf("The subscribe with id = %d is not found", idInt),
			"",
		)
		errDto.Write(w)
		return
	} else if err != nil {
		errDto = models.NewFullExceptionDto(
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to delete the subscribe with id = %d", idInt),
			err.Error(),
		)
		errDto.Write(w)
		return
	}

	// result
	w.WriteHeader(http.StatusNoContent)
	w.Header().Del("Content-Type")
}
