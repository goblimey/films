package people

import (
	personModel "github.com/goblimey/films/models/person"
)

// The ListForm holds view data including a list of people.  It's approximately equivalent 
// to a Struts form bean.
type ListForm interface {
	// People returns the list of Person objects from the form
	People() []personModel.Person
	// Notice gets the notice.
	Notice() string
	// ErrorMessage gets the general error message.
	ErrorMessage() string
	// SetPeople sets the list of Persons in the form.
	SetPeople([]personModel.Person)
	// SetNotice sets the notice.
	SetNotice(notice string)
	//SetErrorMessage sets the error message.
	SetErrorMessage(errorMessage string)
}