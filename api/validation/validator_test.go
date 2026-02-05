package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=18"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "Valid struct",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			input: TestStruct{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			input: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "Underage",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
