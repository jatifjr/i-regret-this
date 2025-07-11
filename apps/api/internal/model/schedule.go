package model

import (
	"time"
)

type Schedule struct {
	ID        int64     `json:"id"`
	PlotID    int64     `json:"plot_id"`
	DateTime  time.Time `json:"date_time"`
	Location  string    `json:"location"`
	Quota     int       `json:"quota"`
	Available int       `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListSchedule []Schedule

type PaginatedResponse struct {
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
	TotalItems int          `json:"total_items"`
	Data       ListSchedule `json:"data"`
}

type CreateSchedule struct {
	DateTime time.Time `json:"date_time" validate:"required,notpastdate"`
	Location string    `json:"location" validate:"required"`
	Quota    int       `json:"quota" validate:"required"`
}

type UpdateSchedule struct {
	DateTime time.Time `json:"date_time" validate:"omitempty,notpastdate"`
	Location string    `json:"location" validate:"omitempty"`
	Quota    int       `json:"quota" validate:"omitempty"`
}
