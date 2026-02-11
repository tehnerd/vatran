package models

import (
	"net/http"

	"github.com/tehnerd/vatran/go/katran"
)

// APIErrorCode represents standardized API error codes.
type APIErrorCode string

const (
	// CodeInvalidRequest indicates the request was malformed or invalid.
	CodeInvalidRequest APIErrorCode = "INVALID_REQUEST"
	// CodeNotFound indicates the requested resource was not found.
	CodeNotFound APIErrorCode = "NOT_FOUND"
	// CodeAlreadyExists indicates the resource already exists.
	CodeAlreadyExists APIErrorCode = "ALREADY_EXISTS"
	// CodeSpaceExhausted indicates maximum capacity was reached.
	CodeSpaceExhausted APIErrorCode = "SPACE_EXHAUSTED"
	// CodeBPFFailed indicates a BPF operation failed.
	CodeBPFFailed APIErrorCode = "BPF_FAILED"
	// CodeFeatureDisabled indicates the requested feature is not enabled.
	CodeFeatureDisabled APIErrorCode = "FEATURE_DISABLED"
	// CodeInternalError indicates an internal server error.
	CodeInternalError APIErrorCode = "INTERNAL_ERROR"
	// CodeLBNotInitialized indicates the load balancer is not initialized.
	CodeLBNotInitialized APIErrorCode = "LB_NOT_INITIALIZED"
	// CodeUnauthorized indicates the request lacks valid authentication.
	CodeUnauthorized APIErrorCode = "UNAUTHORIZED"
	// CodeForbidden indicates the request is not allowed.
	CodeForbidden APIErrorCode = "FORBIDDEN"
	// CodeHCServiceUnavailable indicates the healthcheck service is unavailable.
	CodeHCServiceUnavailable APIErrorCode = "HC_SERVICE_UNAVAILABLE"
)

// APIError represents an error response from the API.
type APIError struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
}

// MapKatranError maps a katran error to an APIError and HTTP status code.
//
// Parameters:
//   - err: The error to map.
//
// Returns the HTTP status code and APIError.
func MapKatranError(err error) (int, *APIError) {
	if err == nil {
		return http.StatusOK, nil
	}

	ke, ok := err.(*katran.KatranError)
	if !ok {
		return http.StatusInternalServerError, &APIError{
			Code:    CodeInternalError,
			Message: err.Error(),
		}
	}

	switch ke.Code {
	case katran.ErrInvalidArgument:
		return http.StatusBadRequest, &APIError{
			Code:    CodeInvalidRequest,
			Message: ke.Message,
		}
	case katran.ErrNotFound:
		return http.StatusNotFound, &APIError{
			Code:    CodeNotFound,
			Message: ke.Message,
		}
	case katran.ErrAlreadyExists:
		return http.StatusConflict, &APIError{
			Code:    CodeAlreadyExists,
			Message: ke.Message,
		}
	case katran.ErrSpaceExhausted:
		return http.StatusInsufficientStorage, &APIError{
			Code:    CodeSpaceExhausted,
			Message: ke.Message,
		}
	case katran.ErrBPFFailed:
		return http.StatusInternalServerError, &APIError{
			Code:    CodeBPFFailed,
			Message: ke.Message,
		}
	case katran.ErrFeatureDisabled:
		return http.StatusNotImplemented, &APIError{
			Code:    CodeFeatureDisabled,
			Message: ke.Message,
		}
	default:
		return http.StatusInternalServerError, &APIError{
			Code:    CodeInternalError,
			Message: ke.Message,
		}
	}
}

// NewAPIError creates a new APIError.
//
// Parameters:
//   - code: The error code.
//   - message: The error message.
//
// Returns a new APIError.
func NewAPIError(code APIErrorCode, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewInvalidRequestError creates a new invalid request error.
//
// Parameters:
//   - message: The error message.
//
// Returns a new APIError with INVALID_REQUEST code.
func NewInvalidRequestError(message string) *APIError {
	return NewAPIError(CodeInvalidRequest, message)
}

// NewNotFoundError creates a new not found error.
//
// Parameters:
//   - message: The error message.
//
// Returns a new APIError with NOT_FOUND code.
func NewNotFoundError(message string) *APIError {
	return NewAPIError(CodeNotFound, message)
}

// NewInternalError creates a new internal error.
//
// Parameters:
//   - message: The error message.
//
// Returns a new APIError with INTERNAL_ERROR code.
func NewInternalError(message string) *APIError {
	return NewAPIError(CodeInternalError, message)
}

// NewLBNotInitializedError creates a new LB not initialized error.
//
// Returns a new APIError with LB_NOT_INITIALIZED code.
func NewLBNotInitializedError() *APIError {
	return NewAPIError(CodeLBNotInitialized, "load balancer is not initialized")
}

// NewFeatureDisabledError creates a new feature disabled error.
//
// Parameters:
//   - message: The error message.
//
// Returns a new APIError with FEATURE_DISABLED code.
func NewFeatureDisabledError(message string) *APIError {
	return NewAPIError(CodeFeatureDisabled, message)
}

// NewHCServiceUnavailableError creates a new HC service unavailable error.
//
// Parameters:
//   - message: The error message.
//
// Returns a new APIError with HC_SERVICE_UNAVAILABLE code.
func NewHCServiceUnavailableError(message string) *APIError {
	return NewAPIError(CodeHCServiceUnavailable, message)
}
