package generate

import (
	"testing"
)

func TestConvertIDNaming(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ID_only",
			input:    "ID",
			expected: "Id",
		},
		{
			name:     "user_id",
			input:    "UserID",
			expected: "UserId",
		},
		{
			name:     "product_id",
			input:    "ProductID",
			expected: "ProductId",
		},
		{
			name:     "order_id",
			input:    "OrderID",
			expected: "OrderId",
		},
		{
			name:     "no_id",
			input:    "UserName",
			expected: "UserName",
		},
		{
			name:     "id_at_beginning",
			input:    "IDGenerator",
			expected: "IDGenerator",
		},
		{
			name:     "id_in_middle",
			input:    "UserIDGenerator",
			expected: "UserIDGenerator",
		},
		{
			name:     "lowercase_id",
			input:    "userId",
			expected: "userId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertIDNaming(tt.input)
			if result != tt.expected {
				t.Errorf("convertIDNaming(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
