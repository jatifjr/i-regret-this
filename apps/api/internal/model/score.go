package model

import (
	"time"
)

// Base model
type Score struct {
	ID                         int64     `json:"id"`
	StudentID                  int64     `json:"student_id"`
	TestPlotID                 int64     `json:"test_plot_id"` // format: year+month+date+order = 20250501001
	ListeningComprehension     int       `json:"listening_comprehension" validate:"min=0,max=68"`
	StructureWrittenExpression int       `json:"structure_written_expression" validate:"min=0,max=68"`
	ReadingComprehension       int       `json:"reading_comprehension" validate:"min=0,max=67"`
	TotalScore                 int       `json:"total_score" validate:"min=0,max=677"` // Calculated field: (listening + structure + reading) * 10/3
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// Create model
type CreateScore struct {
	StudentID                  int64 `json:"student_id" validate:"required"`
	TestPlotID                 int64 `json:"test_plot_id" validate:"required"`
	ListeningComprehension     int   `json:"listening_comprehension" validate:"required,min=0,max=68"`
	StructureWrittenExpression int   `json:"structure_written_expression" validate:"required,min=0,max=68"`
	ReadingComprehension       int   `json:"reading_comprehension" validate:"required,min=0,max=67"`
}

// Update model
type UpdateScore struct {
	ListeningComprehension     int `json:"listening_comprehension,omitempty" validate:"omitempty,min=0,max=68"`
	StructureWrittenExpression int `json:"structure_written_expression,omitempty" validate:"omitempty,min=0,max=68"`
	ReadingComprehension       int `json:"reading_comprehension,omitempty" validate:"omitempty,min=0,max=67"`
}
