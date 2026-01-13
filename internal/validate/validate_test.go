package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"US uppercase", "US", false},
		{"us lowercase", "us", false},
		{"EU uppercase", "EU", false},
		{"eu lowercase", "eu", false},
		{"Us mixed", "Us", false},
		{"invalid XX", "XX", true},
		{"empty", "", true},
		{"USA invalid", "USA", true},
		{"Europe invalid", "Europe", true},
		{"whitespace", " US ", true}, // strict, no trimming
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Region(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid numeric", "12345", false},
		{"single digit", "1", false},
		{"large number", "9999999999", false},
		{"alphabetic", "abc", true},
		{"empty", "", true},
		{"decimal", "12.34", true},
		{"negative", "-1", true},
		{"zero", "0", true},
		{"mixed alphanumeric", "123abc", true},
		{"whitespace", " 123 ", true}, // strict, no trimming
		{"leading zeros", "00123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AccountID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantWarning bool
		wantErr     bool
	}{
		{"valid NRAK key", "NRAK-ABCDEFGHIJ1234567890", false, false},
		{"valid NRAK short", "NRAK-1234567890AB", false, false},
		{"NRAI key warning", "NRAI-ABCDEFGHIJ1234567890", true, false},
		{"no prefix warning", "ABCDEFGHIJ1234567890WXYZ", true, false},
		{"too short error", "NRAK-short", false, true},
		{"way too short", "abc", false, true},
		{"empty error", "", false, true},
		{"exactly 16 chars", "1234567890123456", true, false}, // no prefix, warning
		{"15 chars error", "123456789012345", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning, err := APIKey(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, warning, "should not return warning when error")
			} else {
				assert.NoError(t, err)
				if tt.wantWarning {
					assert.NotEmpty(t, warning)
				} else {
					assert.Empty(t, warning)
				}
			}
		})
	}
}

func TestAPIKey_WarningMessage(t *testing.T) {
	warning, err := APIKey("NRAI-ABCDEFGHIJ1234567890")
	assert.NoError(t, err)
	assert.Contains(t, warning, "NRAK-")
	assert.Contains(t, warning, "User API keys")
}

func TestAccountID_ErrorMessage(t *testing.T) {
	err := AccountID("abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "abc")
	assert.Contains(t, err.Error(), "numeric")
}

func TestRegion_ErrorMessage(t *testing.T) {
	err := Region("XX")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "XX")
	assert.Contains(t, err.Error(), "US or EU")
}
