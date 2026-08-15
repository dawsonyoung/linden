package errs

import (
	"errors"
	"fmt"
	"testing"
)

func Test_New_CarriesCodeAndMessage(t *testing.T) {
	err := New(NotFound, "model missing")

	if got := CodeOf(err); got != NotFound {
		t.Errorf("CodeOf = %q, want %q", got, NotFound)
	}
	if got := err.Error(); got != "model missing" {
		t.Errorf("Error() = %q, want %q", got, "model missing")
	}
}

func Test_Wrap_NilError_ReturnsNil(t *testing.T) {
	if got := Wrap(Internal, "context", nil); got != nil {
		t.Fatalf("Wrap(nil) = %v, want nil", got)
	}
}

func Test_Wrap_PreservesCauseForErrorsIs(t *testing.T) {
	sentinel := errors.New("connection refused")

	err := Wrap(Unavailable, "reach ollama", sentinel)

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is could not find the wrapped cause")
	}
	if got := CodeOf(err); got != Unavailable {
		t.Errorf("CodeOf = %q, want %q", got, Unavailable)
	}
}

func Test_Wrap_MessageIncludesCause(t *testing.T) {
	err := Wrap(Unavailable, "reach ollama", errors.New("connection refused"))

	const want = "reach ollama: connection refused"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func Test_Wrap_EmptyMessage_ReportsCauseOnly(t *testing.T) {
	err := Wrap(Internal, "", errors.New("boom"))

	if got := err.Error(); got != "boom" {
		t.Errorf("Error() = %q, want %q", got, "boom")
	}
}

func Test_CodeOf_UnclassifiedError_IsInternal(t *testing.T) {
	if got := CodeOf(errors.New("plain")); got != Internal {
		t.Errorf("CodeOf = %q, want %q", got, Internal)
	}
}

// A zero-value Code must not escape as an empty string; the spec states an
// error carrying no code is treated as internal.
func Test_CodeOf_EmptyCode_IsInternal(t *testing.T) {
	for _, err := range []error{
		New("", "no code supplied"),
		Wrap("", "no code supplied", errors.New("cause")),
		&Error{},
	} {
		if got := CodeOf(err); got != Internal {
			t.Errorf("CodeOf(%v) = %q, want %q", err, got, Internal)
		}
	}
}

func Test_Error_NoMessageNoCause_ReportsCode(t *testing.T) {
	err := New(Internal, "")

	if got := err.Error(); got != string(Internal) {
		t.Errorf("Error() = %q, want %q", got, string(Internal))
	}
}

func Test_Unwrap_ReturnsTheOriginalCause(t *testing.T) {
	cause := errors.New("connection refused")

	err := Wrap(Unavailable, "reach ollama", cause)

	var target *Error
	if !errors.As(err, &target) {
		t.Fatal("errors.As could not extract *Error")
	}
	if got := target.Unwrap(); got != cause {
		t.Errorf("Unwrap() = %v, want the original cause", got)
	}
}

func Test_CodeOf_Nil_IsEmpty(t *testing.T) {
	if got := CodeOf(nil); got != "" {
		t.Errorf("CodeOf(nil) = %q, want empty", got)
	}
}

func Test_CodeOf_FindsCodeThroughForeignWrapping(t *testing.T) {
	inner := New(DeadlineExceeded, "provider timed out")
	outer := fmt.Errorf("orchestrating chat: %w", inner)

	if got := CodeOf(outer); got != DeadlineExceeded {
		t.Errorf("CodeOf = %q, want %q", got, DeadlineExceeded)
	}
}

func Test_CodeOf_NestedCodes_ReturnsOutermost(t *testing.T) {
	inner := New(Unavailable, "transport failed")
	outer := Wrap(Internal, "unexpected", inner)

	// The outermost classification wins: a layer that reclassifies a failure is
	// making a deliberate statement about it.
	if got := CodeOf(outer); got != Internal {
		t.Errorf("CodeOf = %q, want %q", got, Internal)
	}
}

func Test_Is_MatchesCarriedCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code Code
		want bool
	}{
		{"matching code", New(NotFound, "x"), NotFound, true},
		{"different code", New(NotFound, "x"), Conflict, false},
		{"wrapped match", Wrap(Unavailable, "y", errors.New("z")), Unavailable, true},
		{"unclassified is internal", errors.New("plain"), Internal, true},
		{"nil error", nil, Internal, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.code); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ErrorsAs_ExposesCodeField(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", New(ResourceExhausted, "too many"))

	var target *Error
	if !errors.As(err, &target) {
		t.Fatal("errors.As could not extract *Error")
	}
	if target.Code != ResourceExhausted {
		t.Errorf("Code = %q, want %q", target.Code, ResourceExhausted)
	}
}

// The taxonomy is a published contract; this fails if a constant is renamed or
// its string value drifts from docs/architecture/layer-interface-spec.md.
//
// It cannot detect a tenth constant added without a documentation update. That
// remains a review responsibility under the charter in ADR-0004.
func Test_Taxonomy_MatchesPublishedCodes(t *testing.T) {
	published := map[Code]string{
		InvalidArgument:   "invalid_argument",
		Unauthenticated:   "unauthenticated",
		PermissionDenied:  "permission_denied",
		NotFound:          "not_found",
		Conflict:          "conflict",
		ResourceExhausted: "resource_exhausted",
		Unavailable:       "unavailable",
		DeadlineExceeded:  "deadline_exceeded",
		Internal:          "internal",
	}

	if len(published) != 9 {
		t.Fatalf("taxonomy has %d codes, want 9", len(published))
	}
	for code, want := range published {
		if string(code) != want {
			t.Errorf("code %v serializes as %q, want %q", code, string(code), want)
		}
	}
}
