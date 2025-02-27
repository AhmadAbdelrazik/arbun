package validator

import (
	"fmt"
	"regexp"
)

var (
	EmailRX           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	EgyPhoneNumbersRX = regexp.MustCompile("^(\\+20|0)1[0125]\\d{8}$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, value string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(condition bool, key, value string) {
	if condition == false {
		v.AddError(key, value)
	}
}
func (v *Validator) Error() string {
	result := ""
	for name, error := range v.Errors {
		result += fmt.Sprintf("%s: %s;", name, error)
	}

	return result
}

func (v *Validator) Err() *Validator {
	if len(v.Errors) == 0 {
		return nil
	} else {
		return v
	}
}

func (v *Validator) Matches(word string, regex regexp.Regexp) bool {
	return regex.MatchString(word)
}

func (v *Validator) In(item string, items ...string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false
}

func (v *Validator) Add(vv *Validator) {
	if vv == nil {
		return
	}
	for key, value := range vv.Errors {
		v.AddError(key, value)
	}
}
