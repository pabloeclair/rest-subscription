package models

import (
	"errors"
	"time"
)

type Subscribe struct {
	ID          uint `gorm:"primaryKey"`
	ServiceName string
	Price       int
	UserId      string
	StartDate   time.Time
	EndDate     *time.Time
}

func (s *Subscribe) ToDto() *SubscribeDto {
	return &SubscribeDto{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       &s.Price,
		UserId:      s.UserId,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
	}
}

type SubscribeDto struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       *int       `json:"price"`
	UserId      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func (s *SubscribeDto) Validate() error {
	if s.ServiceName == "" || s.Price == nil || s.UserId == "" || s.StartDate.IsZero() {
		return errors.New("The fields 'service_name', 'price', 'user_id' and 'start_date' are required")
	}
	return s.ValidateTime()
}

func (s *SubscribeDto) ValidateTime() error {
	if s.EndDate != nil && s.StartDate.After(*s.EndDate) {
		return errors.New("The field 'end_time' should be after the 'start_time'")
	}
	return nil
}

func (s *SubscribeDto) ToDatabase() *Subscribe {
	return &Subscribe{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       *s.Price,
		UserId:      s.UserId,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
	}
}
