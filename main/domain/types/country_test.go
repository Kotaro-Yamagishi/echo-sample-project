package types

import (
	"testing"
)

func TestValidateCountryName_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "正常系: 有効な国名",
			input:       "Japan",
			expectError: false,
		},
		{
			name:        "エラー系: 空文字",
			input:       "",
			expectError: true,
			errorMsg:    "country name is required",
		},
		{
			name:        "正常系: 長い国名",
			input:       "ThisIsAVeryLongCountryNameThatExceedsTheMaximumLengthLimit",
			expectError: false,
		},
		{
			name:        "正常系: 50文字の国名",
			input:       "ThisIsAVeryLongCountryNameThatIsExactlyFiftyCharsLong",
			expectError: false,
		},
		{
			name:        "正常系: 空白文字を含む国名",
			input:       "   ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCountryName(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
