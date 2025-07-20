package models

import (
	"encoding/json"
	"net/http"
)

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
	var errDto FullExceptionDto
	if err := json.Unmarshal(b, &errDto); err != nil {
		return 0, err
	}
	if errDto.ErrorMessage != "" {
		lrw.StatusMessage += ": " + errDto.ErrorMessage + ": " + errDto.FullErrorMessage
	}
	errDto.FullErrorMessage = ""
	b, err := json.Marshal(&errDto)
	if err != nil {
		return 0, err
	}
	return lrw.ResponseWriter.Write(b)
}
