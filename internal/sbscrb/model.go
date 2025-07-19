package sbscrb

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SubscribeDto struct {
	gorm.Model
	Id          uint      `gorm:"primaryKey" json:"id"`
	ServiceName string    `gorm:"size:255;not null;default:null" json:"service_name"`
	Price       int       `gorm:"not null;default:null" json:"price"`
	UserId      string    `gorm:"not null;default:null" json:"user_id"`
	StartDate   time.Time `gorm:"not null;default:null" json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type ExceptionDto struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

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
	var errDto ExceptionDto
	if err := json.Unmarshal(b, &errDto); err != nil {
		return 0, err
	}
	if errDto.ErrorMessage != "" {
		lrw.StatusMessage += ": " + errDto.ErrorMessage
	}
	return lrw.ResponseWriter.Write(b)
}
