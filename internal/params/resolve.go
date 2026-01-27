package params

import (
	"fmt"
)

// ResolveParameters merges CLI values with defaults from SQL
// Priority: CLI flags > defaults
// Returns map of param name -> value
func ResolveParameters(paramDefs, cliValues map[string]string) map[string]string {
	result := make(map[string]string)

	// First, add all defaults
	for name, defaultValue := range paramDefs {
		result[name] = defaultValue
	}

	// Then override with CLI values (higher priority)
	for name, cliValue := range cliValues {
		if _, exists := result[name]; exists {
			result[name] = cliValue
		}
	}

	return result
}

// GetMissingRequired finds parameters that have no value (empty string = required)
// Returns list of param names that need user input
func GetMissingRequired(paramDefs, currentValues map[string]string) []string {
	var missing []string

	for name, defaultValue := range paramDefs {
		// If default is empty, it's required
		if defaultValue == "" {
			if value, exists := currentValues[name]; !exists || value == "" {
				missing = append(missing, name)
			}
		}
	}

	return missing
}

// ValidateCLIValues checks if CLI-provided values are valid
// Returns error if a CLI param doesn't exist in the query
func ValidateCLIValues(cliValues, paramDefs map[string]string) error {
	for name := range cliValues {
		if _, exists := paramDefs[name]; !exists {
			return fmt.Errorf("unknown parameter: %s", name)
		}
	}
	return nil
}

// Reserved flags that cannot be used as parameter names
var reservedFlags = map[string]bool{
	"edit":    true,
	"last":    true,
	"l":       true,
	"help":    true,
	"h":       true,
	"version": true,
	"v":       true,
}

// ValidateParamNames checks if parameter names conflict with reserved flags
func ValidateParamNames(paramDefs map[string]string) error {
	for name := range paramDefs {
		if reservedFlags[name] {
			return fmt.Errorf("parameter name '%s' conflicts with reserved flag", name)
		}
	}
	return nil
}
