package job

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrEmptyType   = errors.New("job type is required")
	ErrInvalidData = errors.New("job payload must be a JSON object")
)

type CreateInput struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
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
