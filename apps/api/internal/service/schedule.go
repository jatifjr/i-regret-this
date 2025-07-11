package service

import (
	"context"

	"github.com/jatifjr/app-unw-toefl/apps/api/internal/model"
	"github.com/jatifjr/app-unw-toefl/apps/api/internal/repository"
	"github.com/jatifjr/app-unw-toefl/apps/api/pkg/validator"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, req *model.CreateSchedule) (*model.Schedule, error) {
	if err := validator.New().Struct(req); err != nil {
		return nil, err
	}

	schedule := &model.Schedule{
		DateTime:  req.DateTime,
		Location:  req.Location,
		Quota:     req.Quota,
		Available: req.Quota,
	}

	if err := s.repo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id int64) (*model.Schedule, error) {
	schedule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id int64, req *model.UpdateSchedule) (*model.Schedule, error) {
	if err := validator.New().Struct(req); err != nil {
		return nil, err
	}

	schedule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	schedule.DateTime = req.DateTime
	schedule.Location = req.Location
	schedule.Quota = req.Quota

	if err := s.repo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id int64) error {
	// Get schedule first to check if it exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *ScheduleService) ListSchedules(ctx context.Context, page, pageSize int, sortBy, sortOrder string) (*model.PaginatedResponse, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	schedules, total, err := s.repo.List(ctx, limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Convert []*model.Schedule to []model.Schedule
	scheduleList := make([]model.Schedule, len(schedules))
	for i, s := range schedules {
		if s != nil {
			scheduleList[i] = *s
		}
	}

	// Calculate total pages
	totalPages := (int(total) + pageSize - 1) / pageSize

	return &model.PaginatedResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: int(total),
		Data:       scheduleList,
	}, nil
}
