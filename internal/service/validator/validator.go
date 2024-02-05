package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator checks and validates user data and stores erros in
// FieldErrors map if there any
type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(isRight bool, key, message string) {
	if !isRight {
		v.AddFieldError(key, message)
	}
}

func NotBlank(str string) bool {
	return strings.TrimSpace(str) != ""
}

func MaxChar(str string, n int) bool {
	return utf8.RuneCountInString(str) <= n
}

func MinChar(str string, n int) bool {
	return utf8.RuneCountInString(str) >= n
}

func Matches(str string, rx *regexp.Regexp) bool {
	return rx.MatchString(str)
}
