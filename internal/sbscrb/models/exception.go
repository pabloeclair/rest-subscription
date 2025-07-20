package models

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatih/color"
)

var RedString = color.New(color.FgRed).SprintFunc()

type ExceptionDto struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

type FullExceptionDto struct {
	ExceptionDto
	FullErrorMessage string `json:"full_error_message"`
}

func NewFullExceptionDto(status int, msg string, err string) FullExceptionDto {
	return FullExceptionDto{
		ExceptionDto: ExceptionDto{
			StatusCode:   status,
			ErrorMessage: msg,
		},
		FullErrorMessage: err,
	}
}

func (errDto *FullExceptionDto) Write(w http.ResponseWriter) {
	b, err := json.Marshal(errDto)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println(RedString("ERROR: json marshal: ExceptionDto ", errDto))
		return
	}
	w.WriteHeader(errDto.StatusCode)
	w.Write(b)
}
