package validate

import (
	"fmt"
	"strconv"
	"strings"
)

// Region validates New Relic region (US or EU)
func Region(region string) error {
	upper := strings.ToUpper(region)
	if upper != "US" && upper != "EU" {
		return fmt.Errorf("invalid region %q: must be US or EU", region)
	}
	return nil
}

// AccountID validates account ID is numeric and positive
func AccountID(id string) error {
	if id == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	num, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid account ID %q: must be numeric", id)
	}

	if num <= 0 {
		return fmt.Errorf("invalid account ID %q: must be a positive number", id)
	}

	return nil
}

// APIKey validates API key format
// Returns warning message (not error) for non-standard formats
func APIKey(key string) (warning string, err error) {
	if key == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	// Check minimum length (NRAK- keys are typically 40+ chars)
	if len(key) < 16 {
		return "", fmt.Errorf("API key too short: minimum 16 characters")
	}

	// Check for NRAK- prefix (user keys)
	if !strings.HasPrefix(key, "NRAK-") {
		return "API key does not start with 'NRAK-' (expected for User API keys)", nil
	}

	return "", nil
}
