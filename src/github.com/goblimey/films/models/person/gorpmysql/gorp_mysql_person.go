package gorpmysql

import (
	"fmt"
	"strings"

	personModel "github.com/goblimey/films/models/person"
)

// The GorpMysqlPerson struct implements the Person interface and holds a single row from
// the PEOPLE table, accessed via the GORP library.
//
// The fields must be public for GORP to work and the names must not clash with those of the getters
type GorpMysqlPerson struct {
	IDField       uint64 `db: "id, primarykey, autoincrement"`
	ForenameField string `db: "forename"`
	SurnameField  string `db: "surname"`
}

// Factory functions

// MakePerson creates and returns a new uninitialised Person object
func MakePerson() personModel.Person {
	var GorpMysqlPerson GorpMysqlPerson
	return &GorpMysqlPerson
}

// MakeInitialisedPerson creates and returns a new Person object initialised from
// the arguments
func MakeInitialisedPerson(id uint64, forename string, surname string) personModel.Person {
	person := MakePerson()
	person.SetID(id)
	person.SetForename(forename)
	person.SetSurname(surname)
	return person
}

// Clone creates and returns a new Person object initialised from a source Person.
func Clone(source personModel.Person) personModel.Person {
	return MakeInitialisedPerson(source.ID(), source.Forename(), source.Surname())
}

// Methods to implement the Person interface.

// ID gets the id of the person.
func (p GorpMysqlPerson) ID() uint64 {
	return p.IDField
}

// Forename gets the forename of the person
func (p GorpMysqlPerson) Forename() string {
	return p.ForenameField
}

// Surname gets the surname of the person
func (p GorpMysqlPerson) Surname() string {
	return p.SurnameField
}

// String renders the person as a string
func (p GorpMysqlPerson) String() string {
	return fmt.Sprintf("{%d, %s, %s}", p.IDField, p.ForenameField, p.SurnameField)
}

// SetID sets the person's id to the given value
func (p *GorpMysqlPerson) SetID(id uint64) {
	p.IDField = id
}

// SetForename sets the person's forename to the given value
func (p *GorpMysqlPerson) SetForename(forename string) {
	p.ForenameField = strings.TrimSpace(forename)
}

// SetSurname sets the person's surname to the given value
func (p *GorpMysqlPerson) SetSurname(surname string) {
	p.SurnameField = strings.TrimSpace(surname)
}
