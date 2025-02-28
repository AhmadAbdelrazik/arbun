package validator

import "fmt"

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

func (v *Validator) Add(vv *Validator) {
	if vv == nil {
		return
	}
	for key, value := range vv.Errors {
		v.AddError(key, value)
	}
}
