// Package errs defines the error taxonomy published in
// docs/architecture/layer-interface-spec.md.
//
// It is a shared kernel: standard library only, no I/O, no transport knowledge.
// Mapping a Code to an HTTP status or a client-facing message belongs to the api
// layer. See docs/adr/0004-shared-error-taxonomy.md for the charter.
package errs

import "errors"

// Code identifies a class of failure. The set is closed; adding one requires
// updating the interface specification in the same change.
type Code string

const (
	InvalidArgument   Code = "invalid_argument"
	Unauthenticated   Code = "unauthenticated"
	PermissionDenied  Code = "permission_denied"
	NotFound          Code = "not_found"
	Conflict          Code = "conflict"
	ResourceExhausted Code = "resource_exhausted"
	Unavailable       Code = "unavailable"
	DeadlineExceeded  Code = "deadline_exceeded"
	Internal          Code = "internal"
)

// Error carries a taxonomy code alongside the underlying cause.
//
// Match on the code with Is or errors.As. Two Errors with the same code are not
// equal under errors.Is, which compares identity.
type Error struct {
	Code Code

	// Message is for server-side diagnosis. It reaches logs, so it must contain
	// no user content, and it may name internal detail, so it must never be
	// returned to a client.
	Message string

	cause error
}

// Error returns a diagnostic string for logs. It is not client-facing.
func (e *Error) Error() string {
	switch {
	case e.cause == nil && e.Message == "":
		return string(e.Code)
	case e.cause == nil:
		return e.Message
	case e.Message == "":
		return e.cause.Error()
	default:
		return e.Message + ": " + e.cause.Error()
	}
}

func (e *Error) Unwrap() error { return e.cause }

// New returns an error carrying code.
func New(code Code, message string) error {
	return &Error{Code: code, Message: message}
}

// Wrap annotates err with code, preserving the chain for errors.Is and
// errors.As. It returns nil when err is nil, so a caller can wrap
// unconditionally.
//
// Wrap reclassifies. To add context without claiming a new classification, use
// fmt.Errorf with %w: wrapping in Internal merely to add a note would downgrade
// a precise code and is the most likely way to misreport a failure.
func Wrap(code Code, message string, err error) error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Message: message, cause: err}
}

// CodeOf reports the code carried by the outermost Error in err's chain.
//
// An error with no code is Internal: an unclassified failure is a defect, and
// treating it as anything more specific would misreport it.
func CodeOf(err error) Code {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) && e.Code != "" {
		return e.Code
	}
	return Internal
}

// Is reports whether err carries code.
func Is(err error, code Code) bool {
	return err != nil && CodeOf(err) == code
}
