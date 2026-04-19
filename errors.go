package spidra

import (
	"encoding/json"
	"fmt"
)

// SpidraError is the base error type returned for all API errors.
type SpidraError struct {
	StatusCode int
	Message    string
}

func (e *SpidraError) Error() string {
	return fmt.Sprintf("spidra: %d: %s", e.StatusCode, e.Message)
}

// AuthenticationError is returned when the API key is missing or invalid (401).
type AuthenticationError struct{ *SpidraError }

// InsufficientCreditsError is returned when the account has no remaining credits (403).
type InsufficientCreditsError struct{ *SpidraError }

// RateLimitError is returned when too many requests are sent (429).
type RateLimitError struct{ *SpidraError }

// ServerError is returned when Spidra's servers encounter an error (5xx).
type ServerError struct{ *SpidraError }

func mapError(status int, body []byte) error {
	var payload struct {
		Error string `json:"error"`
	}
	msg := "unknown error"
	if json.Unmarshal(body, &payload) == nil && payload.Error != "" {
		msg = payload.Error
	}

	base := &SpidraError{StatusCode: status, Message: msg}
	switch status {
	case 401:
		return &AuthenticationError{base}
	case 403:
		return &InsufficientCreditsError{base}
	case 429:
		return &RateLimitError{base}
	default:
		if status >= 500 {
			return &ServerError{base}
		}
		return base
	}
}
