package models

import (
	"encoding/json"
	"net/http"
)

var NotificationInternalError string

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	StatusMessage string
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, "OK"}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.StatusMessage = http.StatusText(code)
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	var errFullDto FullExceptionDto
	if err := json.Unmarshal(b, &errFullDto); err != nil {
		return 0, err
	}
	if errFullDto.ErrorMessage != "" {
		lrw.StatusMessage += ": " + errFullDto.ErrorMessage + ": " + errFullDto.FullErrorMessage
	}

	errDto := ExceptionDto{
		StatusCode:   errFullDto.StatusCode,
		ErrorMessage: errFullDto.ErrorMessage,
	}
	if errFullDto.StatusCode == http.StatusInternalServerError {
		errDto.ErrorMessage += ". " + NotificationInternalError
	}

	b, err := json.Marshal(&errDto)
	if err != nil {
		return 0, err
	}
	return lrw.ResponseWriter.Write(b)
}
