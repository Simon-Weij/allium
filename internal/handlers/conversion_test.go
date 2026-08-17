package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		dateString   string
		expectedYear int
	}{
		{
			name:         "should parse valid RFC3339 date",
			dateString:   "2024-01-15T12:00:00Z",
			expectedYear: 2024,
		},
		{
			name:         "should return zero for invalid date",
			dateString:   "not-a-date",
			expectedYear: 0,
		},
		{
			name:         "should return zero for empty date",
			dateString:   "",
			expectedYear: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expectedYear, parseYear(tt.dateString))
		})
	}
}
