package errors

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	// ErrRecordNotFound represents error when record is not found.
	ErrRecordNotFound = NewError(codes.NotFound, "record not found")
	// ErrInternalServerError represents error when internal server error occurs.
	ErrInternalServerError = NewError(codes.Internal, "internal server error")
	// ErrWrongLoginCredentials represents error when login credentials are wrong.
	ErrWrongLoginCredentials = NewError(codes.InvalidArgument, "username atau password salah")
)

// Error represents a data structure for error.
// It implements golang error interface.
type Error struct {
	// Code represents error code.
	Code codes.Code `json:"code"`
	// Message represents error message.
	// This is the message that exposed to the user.
	Message string `json:"message"`
}

// NewError creates an instance of Error.
func NewError(code codes.Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error returns internal message in one string.
func (err *Error) Error() error {
	return fmt.Errorf("%d:%s", err.Code, err.Message)
}

// ParseError parses error message and returns an instance of Error.
func ParseError(err error) *Error {
	if err == nil {
		return nil
	}

	split := strings.Split(err.Error(), ":")

	var strToCode = map[string]codes.Code{
		"0":  codes.OK,
		"1":  codes.Canceled,
		"2":  codes.Unknown,
		"3":  codes.InvalidArgument,
		"4":  codes.DeadlineExceeded,
		"5":  codes.NotFound,
		"6":  codes.AlreadyExists,
		"7":  codes.PermissionDenied,
		"8":  codes.ResourceExhausted,
		"9":  codes.FailedPrecondition,
		"10": codes.Aborted,
		"11": codes.OutOfRange,
		"12": codes.Unimplemented,
		"13": codes.Internal,
		"14": codes.Unavailable,
		"15": codes.DataLoss,
		"16": codes.Unauthenticated,
	}

	return NewError(strToCode[split[0]], split[1])
}
