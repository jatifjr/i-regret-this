package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jatifjr/app-unw-toefl/apps/api/internal/model"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/service"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req model.CreateSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	schedule, err := h.service.CreateSchedule(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "test date cannot be in the past" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	schedule, err := h.service.GetSchedule(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "schedule not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	var req model.UpdateSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	schedule, err := h.service.UpdateSchedule(c.Request.Context(), id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "schedule not found":
			status = http.StatusNotFound
		case "test date cannot be in the past", "invalid schedule: quota must be greater than 0":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	if err := h.service.DeleteSchedule(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "schedule not found" {
			status = http.StatusNotFound
		} else if err.Error() == "schedule not found or has registrations" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ScheduleHandler) ListSchedules(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get sorting parameters
	sortBy := c.DefaultQuery("sort_by", "date_time")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	// Validate sort parameters
	validSortFields := map[string]bool{
		"date_time": true,
		"available": true,
	}
	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortFields[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort_by field. Valid fields are: date, time, available"})
		return
	}
	if !validSortOrders[sortOrder] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort_order. Valid values are: asc, desc"})
		return
	}

	schedules, err := h.service.ListSchedules(c.Request.Context(), page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) RegisterRoutes(router *gin.RouterGroup) {
	schedules := router.Group("/schedules")
	{
		schedules.POST("", h.CreateSchedule)
		schedules.GET("/:id", h.GetSchedule)
		schedules.PUT("/:id", h.UpdateSchedule)
		schedules.DELETE("/:id", h.DeleteSchedule)
		schedules.GET("", h.ListSchedules)
	}
}
