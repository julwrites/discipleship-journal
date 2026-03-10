package models

import (
	"testing"
)

func TestReadingPlanDay_Validate(t *testing.T) {
	tests := []struct {
		name    string
		passage string
		wantErr bool
	}{
		{
			name:    "Valid single passage",
			passage: "Genesis 1",
			wantErr: false,
		},
		{
			name:    "Valid semicolon separated",
			passage: "Genesis 1; Exodus 2",
			wantErr: false,
		},
		{
			name:    "Invalid newline",
			passage: "Genesis 1\nExodus 2",
			wantErr: true,
		},
		{
			name:    "Invalid ampersand",
			passage: "Genesis 1 & Exodus 2",
			wantErr: true,
		},
		{
			name:    "Invalid plus",
			passage: "Genesis 1 + Exodus 2",
			wantErr: true,
		},
		{
			name:    "Empty passage",
			passage: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ReadingPlanDay{
				Passage: tt.passage,
			}
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ReadingPlanDay.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
