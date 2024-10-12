package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Error represents a user in the system.
// @Description Error model for handling errors.
// @Name Error
// @Property Code string true "Error code"
// @Property Message string true "Error message"
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("{code: %s, message: %s}", e.Code, e.Message)
}
func (e *Error) JSON() []byte {
	jsonError, _ := json.Marshal(e)
	return jsonError
}

var (
	ErrInternalError                 = &Error{Code: "INTERNAL_ERROR", Message: "An internal server error occurred"}
	ErrInvalidInput                  = &Error{Code: "INVALID_INPUT", Message: "Input data is invalid"}
	ErrBadRequest                    = &Error{Code: "BAD_REQUEST", Message: "bad request"}
	ErrInternalServerError           = &Error{Code: "INTERNAL_SERVER_ERROR", Message: "internal server error"}
	ErrUnauthorizationHeaderNotFound = &Error{Code: "UNAUTHORIZATION_HEADER_NOT_FOUND", Message: "Authorization header required"}
	ErrInvalidToken                  = &Error{Code: "INVALID_TOKEN", Message: "Invalid token"}
	ErrInvalidAuthHeader             = &Error{Code: "INVALID_AUTH_HEADER", Message: "Invalid Authorization header format. Expected"}
)

var (
	ErrGroupNotFound = &Error{Code: "GROUP_NOT_FOUND", Message: "The specified group could not be found"}
)

// User-Related Errors
var (
	ErrUserNotFound      = &Error{Code: "USER_NOT_FOUND", Message: "The specified user could not be found"}
	ErrUserAlreadyExists = &Error{Code: "USER_ALREADY_EXISTS", Message: "A user with this email or username already exists"}
	ErrInvalidCredential = &Error{Code: "INVALID_CREDENTIAL", Message: "username or password incorrect"}
)

var (
	ErrPaymentAlreadyMade   = &Error{Code: "PAYMENT_ALREADY_MADE", Message: "Payment has already been made by this user"}
	ErrPaymentFailed        = &Error{Code: "PAYMENT_FAILED", Message: "Failed to update payment status"}
	ErrBillCompletionFailed = &Error{Code: "BILL_COMPLETION_FAILED", Message: "Failed to mark the bill as completed"}
	ErrGroupUpdateFailed    = &Error{Code: "GROUP_UPDATE_FAILED", Message: "Failed to update group information"}
)

// Validation error functions
func ErrRequired(t any) error {
	return &Error{Code: "VALIDATION_REQUIRED", Message: fmt.Sprintf("%s is required", reflect.TypeOf(t).Name())}
}

func ErrBadInput(t any) error {
	return &Error{Code: "BAD_INPUT", Message: fmt.Sprintf("Invalid input for %s", reflect.TypeOf(t).Name())}
}

func ErrInvalid(s string) error {
	return &Error{Code: "INVALID_REQUEST", Message: s}
}

func ErrUsersNotFound(email []string) error {

	return &Error{Code: "USERS_NOT_FOUND", Message: fmt.Sprintf("Users not found with email : %v", strings.Join(email, ", "))}
}

func ErrValidationFailed(s string) error {
	return &Error{Code: "VALIDATION_FAILED", Message: s}
}
func ErrInvalidQueryParameter(s string) error {
	return &Error{Code: "INVALID_QUERY_PARAMETER", Message: s}
}
