package person

import (
	"fmt"
)

// ConcretePerson represents a person and satisfies the Person interface.
type ConcretePerson struct {
	id       uint64
	forename string
	surname  string
}

// Define the factory functions.

// MakePerson creates and returns a new uninitialised Person object
func MakePerson() Person {
	var concretePerson ConcretePerson
	return &concretePerson
}

// MakeInitialisedPerson creates and returns a new Person object initialised from
// the arguments
func MakeInitialisedPerson(id uint64, forename string, surname string) Person {
	person := MakePerson()
	person.SetID(id)
	person.SetForename(forename)
	person.SetSurname(surname)
	return person
}

// Clone creates and returns a new Person object initialised from a source Person.
func Clone(source Person) Person {
	return MakeInitialisedPerson(source.ID(), source.Forename(), source.Surname())
}

// Define the getters.

// ID() gets the id of the person.
func (cp ConcretePerson) ID() uint64 {
	return cp.id
}

//Forename gets the forename of the person.
func (cp ConcretePerson) Forename() string {
	return cp.forename
}

// Surname gets the surname of the person.
func (cp ConcretePerson) Surname() string {
	return cp.surname
}

// String gets the person as a String.
func (cp ConcretePerson) String() string {
	return fmt.Sprintf("ConcretePerson={id=%d, forename=%s,surname=%s}",
		cp.id,
		cp.surname,
		cp.forename)
}

// Define the setters.

// SetID sets the id to the given value.
func (cp *ConcretePerson) SetID(id uint64) {
	cp.id = id
}

// SetForename sets the forename of the person.
func (cp *ConcretePerson) SetForename(forename string) {
	cp.forename = forename
}

// SetSurname sets the surname of the person.
func (cp *ConcretePerson) SetSurname(surname string) {
	cp.surname = surname
}
