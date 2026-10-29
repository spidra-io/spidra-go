package spidra

import (
	"encoding/json"
	"fmt"
)

// SpidraError is the base error type returned for all API errors.
type SpidraError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *SpidraError) Error() string {
	return fmt.Sprintf("spidra: %d: %s", e.StatusCode, e.Message)
}

// AuthenticationError is returned when the API key is missing or invalid (401).
type AuthenticationError struct{ *SpidraError }

// PaymentRequiredError is returned when payment is required (402).
type PaymentRequiredError struct{ *SpidraError }

// ForbiddenError is returned when access is forbidden or credits are exhausted (403).
type ForbiddenError struct{ *SpidraError }

// InsufficientCreditsError is returned when the account has no remaining credits.
type InsufficientCreditsError struct{ *SpidraError }

// NotFoundError is returned when the requested job, preset, or resource does not exist (404).
type NotFoundError struct{ *SpidraError }

// ValidationError is returned when request parameters or body fail validation (422).
type ValidationError struct {
	*SpidraError
	Errors []string
}

// RateLimitError is returned when too many requests are sent (429).
type RateLimitError struct{ *SpidraError }

// ServerError is returned when Spidra's servers encounter an unexpected error (5xx).
type ServerError struct{ *SpidraError }

// JobFailedError is returned when an asynchronous job fails or is cancelled.
type JobFailedError struct {
	JobID   string
	Status  string
	Message string
}

func (e *JobFailedError) Error() string {
	return fmt.Sprintf("spidra: job %s %s: %s", e.JobID, e.Status, e.Message)
}

// TimeoutError is returned when an operation or job poll times out.
type TimeoutError struct {
	JobID   string
	Timeout string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("spidra: job %s timed out after %s", e.JobID, e.Timeout)
}

func mapError(status int, body []byte) error {
	var payload struct {
		Message string   `json:"message"`
		Error   any      `json:"error"`
		Errors  []string `json:"errors"`
	}

	msg := fmt.Sprintf("HTTP %d", status)
	var errDetails []string

	if json.Unmarshal(body, &payload) == nil {
		if payload.Message != "" {
			msg = payload.Message
		} else if str, ok := payload.Error.(string); ok && str != "" {
			msg = str
		} else if m, ok := payload.Error.(map[string]any); ok {
			if em, ok := m["message"].(string); ok && em != "" {
				msg = em
			}
		}
		errDetails = payload.Errors
	}

	base := &SpidraError{
		StatusCode: status,
		Message:    msg,
		Body:       string(body),
	}

	switch status {
	case 401:
		return &AuthenticationError{base}
	case 402:
		return &PaymentRequiredError{base}
	case 403:
		return &ForbiddenError{base}
	case 404:
		return &NotFoundError{base}
	case 422:
		return &ValidationError{
			SpidraError: base,
			Errors:      errDetails,
		}
	case 429:
		return &RateLimitError{base}
	default:
		if status >= 500 {
			return &ServerError{base}
		}
		return base
	}
}
