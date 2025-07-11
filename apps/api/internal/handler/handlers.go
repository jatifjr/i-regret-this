package handler

import (
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/service"
)

// Handler contains all handlers for the application
type Handler struct {
	Schedule *ScheduleHandler
}

// NewHandler creates a new Handler instance
func NewHandler(scheduleService *service.ScheduleService) *Handler {
	return &Handler{
		Schedule: NewScheduleHandler(scheduleService),
	}
}
