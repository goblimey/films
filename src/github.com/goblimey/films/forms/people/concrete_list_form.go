package people

import (
	personModel "github.com/goblimey/films/models/person"
)

// The ConcreteListForm satisfies the ListForm interface and holds view data 
// including a list of people.  It's approximately equivalent
// to a Struts form bean.
type ConcreteListForm struct {
	people       []personModel.Person
	notice       string
	errorMessage string
}

// People returns the list of Person objects from the form
func (clf *ConcreteListForm) People() []personModel.Person {
	return clf.people
}

// Notice gets the notice.
func (clf *ConcreteListForm) Notice() string {
	return clf.notice
}

// ErrorMessage gets the general error message.
func (clf *ConcreteListForm) ErrorMessage() string {
	return clf.errorMessage
}

// SetPeople sets the list of Persons.
func (clf *ConcreteListForm) SetPeople(people []personModel.Person) {
	clf.people = people
}

// SetNotice sets the notice.
func (clf *ConcreteListForm) SetNotice(notice string) {
	clf.notice = notice
}

// SetErrorMessage sets the error message.
func (clf *ConcreteListForm) SetErrorMessage(errorMessage string) {
	clf.errorMessage = errorMessage
}
