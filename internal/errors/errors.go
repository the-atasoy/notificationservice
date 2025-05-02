package errors

import (
	"fmt"
)

type ErrorType string

const (
	ValidationError ErrorType = "validation"

	RetriableError ErrorType = "retriable"

	ProcessingError ErrorType = "processing"
)

type NotificationError struct {
	Type        ErrorType
	Description string
	OriginalErr error
}

func (e *NotificationError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s error: %s - %v", e.Type, e.Description, e.OriginalErr)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Description)
}

func NewValidationError(description string, err error) *NotificationError {
	return &NotificationError{
		Type:        ValidationError,
		Description: description,
		OriginalErr: err,
	}
}

func NewRetriableError(description string, err error) *NotificationError {
	return &NotificationError{
		Type:        RetriableError,
		Description: description,
		OriginalErr: err,
	}
}

func NewProcessingError(description string, err error) *NotificationError {
	return &NotificationError{
		Type:        ProcessingError,
		Description: description,
		OriginalErr: err,
	}
}

func IsValidationError(err error) bool {
	if notifErr, ok := err.(*NotificationError); ok {
		return notifErr.Type == ValidationError
	}
	return false
}

func IsRetriableError(err error) bool {
	if notifErr, ok := err.(*NotificationError); ok {
		return notifErr.Type == RetriableError
	}
	return false
}

func GetErrorType(err error) ErrorType {
	if notifErr, ok := err.(*NotificationError); ok {
		return notifErr.Type
	}
	return ""
}

func GetErrorDescription(err error) string {
	if notifErr, ok := err.(*NotificationError); ok {
		return notifErr.Description
	}
	return err.Error()
}