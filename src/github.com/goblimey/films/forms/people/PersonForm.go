package people

import (
	personModel "github.com/goblimey/films/models/person"
)

// PersonForm holds view data about a Person.  It's used as a data transfer object (DTO)
// in particular for use with views that handle a Person.  (It's approximately equivalent to
// a Struts form bean.)  It contains a Person; a validator function that validates the data
// in the Person and sets the various error messages; a general error message (for errors not
// associated with an individual field of the Person), a notice (for announcement that are
// not about errors) and a set of error messages about individual fields of the Person.  It
// offers getters and setters for the various attributes that it supports.
type PersonForm interface {
	// Person gets the Person embedded in the form.
	Person() personModel.Person
	// Notice gets the notice.
	Notice() string
	// ErrorMessage gets the general error message.
	ErrorMessage() string
	// FieldErrors returns all the field errors as a map.
	FieldErrors() map[string]string
	// ErrorForField returns the error message about a field (may be an empty string).
	ErrorForField(key string) string
	// String returns a string version of the PersonForm.
	String() string
	// SetPerson sets the Person in the form.
	SetPerson(person personModel.Person)
	// SetNotice sets the notice.
	SetNotice(notice string)
	//SetErrorMessage sets the general error message.
	SetErrorMessage(errorMessage string)
	// SetErrorMessageForField sets the error message for a named field
	SetErrorMessageForField(fieldname, errormessage string)
	// Validate validates the data in the Person and sets the various error messages.
	// It returns true if the data is valid, false if there are errors.
	Validate() bool
}
