package errs

import (
	"encoding/json"

	"github.com/donbarrigon/forge/lang"
	"github.com/donbarrigon/forge/str"
)

type ValidationError struct {
	Messages     map[string][]string
	Placeholders map[string][]str.Placeholder
	IsEntity     bool
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Messages:     map[string][]string{},
		Placeholders: map[string][]str.Placeholder{},
	}
}

func (self *ValidationError) Append(field string, message string, ph str.Placeholder) {
	ph.Append("field", field)
	self.Messages[field] = append(self.Messages[field], message)
	self.Placeholders[field] = append(self.Placeholders[field], ph)
}

func (self *ValidationError) AppendM(field string, message string) {
	self.Messages[field] = append(self.Messages[field], message)
	self.Placeholders[field] = append(self.Placeholders[field], str.NewPlaceholder("field", field))
}

func (self *ValidationError) HasErrors() bool {
	if len(self.Messages) > 0 {
		return true
	}
	return false
}

// ================================================================
// Funciones para la interfaz de Errror
// ================================================================

func (self *ValidationError) Error() string {
	b, err := json.Marshal(self.Messages)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func (self *ValidationError) Errors() *Error {
	data := map[string][]string{}
	if len(self.Messages) > 0 {
		for key, messages := range self.Messages {
			m := make([]string, len(messages))
			for i, message := range messages {
				m[i] = self.Placeholders[key][i].Replace(message)
			}
			data[key] = m
		}
		return &Error{
			Status:  UNPROCESSABLE_ENTITY,
			Message: UNPROCESSABLE_ENTITY_MSG,
			Name:    "Validation Error",
			Cause:   "",
			Stack:   "",
			Data:    data,
		}
	}
	return nil
}

func (self *ValidationError) ErrorsT(locale string) *Error {
	data := map[string][]string{}
	if len(self.Messages) > 0 {
		for key, messages := range self.Messages {
			m := make([]string, len(messages))
			for i, message := range messages {
				m[i] = lang.T(locale, message, self.Placeholders[key][i])
			}
			data[key] = m
		}
		return &Error{
			Status:  UNPROCESSABLE_ENTITY,
			Message: lang.T(locale, UNPROCESSABLE_ENTITY_MSG, nil),
			Name:    "Validation Error",
			Cause:   "",
			Stack:   "",
			Data:    data,
		}
	}
	return nil
}
