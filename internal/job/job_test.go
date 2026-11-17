package job

import (
	"errors"
	"testing"
)

func TestCreateInputValidate(t *testing.T) {
	tests := []struct {
		name string
		in   CreateInput
		want error
	}{
		{
			name: "valid object",
			in:   CreateInput{Type: "report", Payload: []byte(`{"format":"csv"}`)},
		},
		{
			name: "empty type",
			in:   CreateInput{Payload: []byte(`{}`)},
			want: ErrEmptyType,
		},
		{
			name: "missing payload",
			in:   CreateInput{Type: "report"},
			want: ErrInvalidData,
		},
		{
			name: "array payload",
			in:   CreateInput{Type: "report", Payload: []byte(`[]`)},
			want: ErrInvalidData,
		},
		{
			name: "null payload",
			in:   CreateInput{Type: "report", Payload: []byte(`null`)},
			want: ErrInvalidData,
		},
		{
			name: "invalid payload",
			in:   CreateInput{Type: "report", Payload: []byte(`{`)},
			want: ErrInvalidData,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			if !errors.Is(err, tc.want) {
				t.Fatalf("Validate() error = %v, want %v", err, tc.want)
			}
		})
	}
}
