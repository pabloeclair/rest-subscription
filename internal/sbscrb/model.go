package sbscrb

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

var RedString = color.New(color.FgRed).SprintFunc()

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
