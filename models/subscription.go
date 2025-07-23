package models

import (
	"time"
)

type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ParseMonthYear(input string) (time.Time, error) {
	return time.Parse("01-2006", input)
}

func FormatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}
