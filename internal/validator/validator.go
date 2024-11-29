package validator

import (
	"regexp"
	"slices"
)

/* Regex for email */
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, msg string) {
	_, exists := v.Errors[key]

	if !exists {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(acceptable bool, key string, msg string) {
	if !acceptable {
		v.AddError(key, msg)
	}
}

func PermittedValue(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}

/* Check whether or not text matches our regex (for emails) */
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
