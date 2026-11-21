package job

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyType   = errors.New("job type is required")
	ErrInvalidData = errors.New("job payload must be a JSON object")
)

type CreateInput struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      Status          `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	Result      json.RawMessage `json:"result,omitempty"`
	LastError   string          `json:"last_error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	FinishedAt  *time.Time      `json:"finished_at,omitempty"`
}

func (in CreateInput) Validate() error {
	if strings.TrimSpace(in.Type) == "" {
		return ErrEmptyType
	}

	var data map[string]json.RawMessage
	if len(in.Payload) == 0 || json.Unmarshal(in.Payload, &data) != nil || data == nil {
		return ErrInvalidData
	}

	return nil
}
