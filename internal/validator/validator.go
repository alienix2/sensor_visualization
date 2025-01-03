package validator

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = message
	}
}
