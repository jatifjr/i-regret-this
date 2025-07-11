package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jatifjr/app-unw-toefl/apps/api/internal/model"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *model.Schedule) error {
	query := `
		INSERT INTO schedules (
			date_time, location, quota
		) VALUES (
			$1, $2, $3
		) RETURNING id, plot_id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		schedule.DateTime,
		schedule.Location,
		schedule.Quota,
	).Scan(
		&schedule.ID,
		&schedule.PlotID,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id int64) (*model.Schedule, error) {
	query := `
		SELECT id, plot_id, date_time, location, 
		       quota, available, created_at, updated_at
		FROM schedules
		WHERE id = $1
	`

	schedule := &model.Schedule{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.PlotID,
		&schedule.DateTime,
		&schedule.Location,
		&schedule.Quota,
		&schedule.Available,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("schedule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return schedule, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *model.Schedule) error {
	query := `
		UPDATE schedules
		SET date_time = $1, location = $2,
		    quota = $3, available = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND available >= 0
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		schedule.DateTime,
		schedule.Location,
		schedule.Quota,
		schedule.Available,
		schedule.ID,
	).Scan(&schedule.UpdatedAt)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("schedule not found or invalid quota")
	}
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM schedules WHERE id = $1 AND available = quota`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("schedule not found or has registrations")
	}

	return nil
}

func (r *ScheduleRepository) List(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]*model.Schedule, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM schedules WHERE date_time >= CURRENT_TIMESTAMP`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Validate and sanitize sort parameters
	validSortFields := map[string]string{
		"date_time": "date_time",
		"available": "available",
	}
	validSortOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	sortField, ok := validSortFields[sortBy]
	if !ok {
		sortField = "date_time" // default sort field
	}
	sortDirection, ok := validSortOrders[sortOrder]
	if !ok {
		sortDirection = "ASC" // default sort direction
	}

	// Build the ORDER BY clause using a CASE statement to prevent SQL injection
	query := `
		SELECT id, plot_id, date_time, location,
		       quota, available, created_at, updated_at
		FROM schedules
		WHERE date_time >= CURRENT_TIMESTAMP
		ORDER BY 
			CASE $3
				WHEN 'date_time' THEN date_time::text
				WHEN 'available' THEN available::text
				ELSE date_time::text
			END ` + sortDirection + `
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset, sortField)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*model.Schedule
	for rows.Next() {
		schedule := &model.Schedule{}
		err := rows.Scan(
			&schedule.ID,
			&schedule.PlotID,
			&schedule.DateTime,
			&schedule.Location,
			&schedule.Quota,
			&schedule.Available,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating schedules: %w", err)
	}

	return schedules, total, nil
}
