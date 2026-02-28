package check

import (
	"fmt"
	"strings"

	"github.com/LazyBachelor/LazyPM/internal/models"
)

type Check = models.Check
type ValidationFeedback = models.ValidationFeedback

func NewCheck(message string, valid bool) Check {
	return Check{
		Message: message,
		Valid:   valid,
	}
}

type Expector struct {
	ValidationFeedback
}

func NewExpector() *Expector {
	return &Expector{
		ValidationFeedback: ValidationFeedback{
			Success: false,
			Checks:  []Check{},
		},
	}
}

func (e *Expector) Complete() ValidationFeedback {
	e.Success = len(e.Errors()) == 0
	return e.ValidationFeedback
}

func (e *Expector) Fail(message string) ValidationFeedback {
	e.Checks = append(e.Checks, NewCheck(message, false))
	return e.ValidationFeedback
}

func (e *Expector) Errors() []error {
	var errors []error
	for _, check := range e.Checks {
		if !check.Valid {
			errors = append(errors, fmt.Errorf("%s", check.Message))
		}
	}
	return errors
}

func (e *Expector) Assert(condition bool, message string) *Expector {
	check := NewCheck(message, condition)
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) Nil(value any, message string) *Expector {
	check := NewCheck(message, value == nil)
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) NotNil(value any, message string) *Expector {
	check := NewCheck(message, value != nil)
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) Contains(s, substr, message string) *Expector {
	check := NewCheck(message, strings.Contains(s, substr))
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) NotContains(s, substr, message string) *Expector {
	check := NewCheck(message, !strings.Contains(s, substr))
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) NotEmptyString(s string, message string) *Expector {
	check := NewCheck(message, s != "")
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) EmptySlice(s []any, message string) *Expector {
	check := NewCheck(message, len(s) == 0)
	e.Checks = append(e.Checks, check)
	return e
}

func (e *Expector) NotEmptySlice(s []any, message string) *Expector {
	check := NewCheck(message, len(s) > 0)
	e.Checks = append(e.Checks, check)
	return e
}
