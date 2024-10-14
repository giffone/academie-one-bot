package domain_test

import (
	"ac_bot/internal/domain"
	"errors"
	"strconv"
	"testing"
)

var (
	// same error
	ErrSame *domain.CliMsg // true
	// different error
	ErrNumError *strconv.NumError // false
	ErrNotSame  *someError        // false
)

func TestCliError(t *testing.T) {
	errExample := &domain.CliMsg{
		Message: "Opps! Error.",
		Err:     errors.New("service: id not equal"),
	}

	var err error = errExample

	// true

	if ok := errors.As(err, &ErrSame); ok {
		if _, ok2 := err.(*domain.CliMsg); !ok2 {
			t.Fatalf("1: error AS: cliMsg method = %t", ok2)
		}
	} else {
		t.Fatalf("2: error AS: equal & equal = %t", ok)
	}

	// false

	if ok := errors.As(err, &ErrNumError); ok {
		t.Fatalf("3: error AS: equal & not equal = %t", ok)
	}

	if ok := errors.As(err, &ErrNotSame); ok {
		t.Fatalf("4: error AS: equal & not equal = %t", ok)
	}
}

type someError struct {
	err error
}

func (e *someError) Error() string {
	return e.err.Error()
}
