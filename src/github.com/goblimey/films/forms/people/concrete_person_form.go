package people

import (
	"fmt"

	personModel "github.com/goblimey/films/models/person"
	"github.com/goblimey/films/utilities"
)

// ConcretePersonForm satisfies the PersonForm interface.
type ConcretePersonForm struct {
	person       personModel.Person
	errorMessage string
	notice       string
	fieldError   map[string]string
}

// Getters

// Person gets the Person embedded in the form.
func (pfd ConcretePersonForm) Person() personModel.Person {
	return pfd.person
}

// Notice gets the notice.
func (pfd ConcretePersonForm) Notice() string {
	return pfd.notice
}

// ErrorMessage gets the general error message.
func (pfd ConcretePersonForm) ErrorMessage() string {
	return pfd.errorMessage
}

// FieldErrors returns all the field errors as a map.
func (pfd ConcretePersonForm) FieldErrors() map[string]string {
	return pfd.fieldError
}

// ErrorForField returns the error message about a field (may be an empty string).
func (pfd ConcretePersonForm) ErrorForField(key string) string {
	if pfd.fieldError == nil {
		// The field error map has not been set up.
		return ""
	}
	return pfd.fieldError[key]
}

// String returns a string version of the PersonForm.
func (pfd ConcretePersonForm) String() string {
	return fmt.Sprintf("ConcretePersonForm={person=%s, notice=%s,errorMessage=%s,fieldError=%s}",
		pfd.person,
		pfd.notice,
		pfd.errorMessage,
		utilities.Map2String(pfd.fieldError))
}

// Setters

// SetPerson sets the Person in the form.
func (pfd *ConcretePersonForm) SetPerson(person personModel.Person) {
	pfd.person = person
}

// SetNotice sets the notice.
func (pfd *ConcretePersonForm) SetNotice(notice string) {
	pfd.notice = notice
}

//SetErrorMessage sets the general error message.
func (pfd *ConcretePersonForm) SetErrorMessage(errorMessage string) {
	pfd.errorMessage = errorMessage
}

// SetErrorMessageForField sets the error message for a named field
func (pfd *ConcretePersonForm) SetErrorMessageForField(fieldname, errormessage string) {
	if pfd.fieldError == nil {
		pfd.fieldError = make(map[string]string)
	}
	pfd.fieldError[fieldname] = errormessage
}

// Validate validates the data in the Person and sets the various error messages.
// It returns true if the data is valid, false if there are errors.
func (pfd *ConcretePersonForm) Validate() bool {
	person := pfd.Person()
	// trim all string items
	person.SetForename(utilities.Trim(person.Forename()))
	person.SetSurname(utilities.Trim(person.Surname()))
	// validate
	valid := true

	if len(person.Forename()) <= 0 {
		pfd.SetErrorMessageForField("Forename", "you must specify the Forename")
		valid = false
	}
	if len(person.Surname()) <= 0 {
		pfd.SetErrorMessageForField("Surname", "you must specify the Surname")
		valid = false
	}
	return valid
}
