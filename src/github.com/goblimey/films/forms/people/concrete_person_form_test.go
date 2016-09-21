package people

import (
	"testing"

	model "github.com/goblimey/films/models/person/gorpmysql"
)

var expectedID uint64 = 42
var expectedForename = "foo"
var expectedSurname = "bar"

// Create a person and a ConcretePersonForm containing it.  Retrieve the person.
func TestUnitCreatePersonFormAndRetrievePerson(t *testing.T) {
	personform := CreatePersonForm(expectedID, expectedForename, expectedSurname)
	if personform.Person().ID() != expectedID {
		t.Errorf("Expected ID to be %d actually %d", expectedID, personform.Person().ID())
	}
	if personform.Person().Forename() != expectedForename {
		t.Errorf("Expected forename to be %s actually %s", expectedForename, personform.Person().Forename())
	}
	if personform.Person().Surname() != expectedSurname {
		t.Errorf("Expected surname to be %s actually %s", expectedSurname, personform.Person().Surname())
	}
}

// Create a personform containing a person with no forename, and validate it.
func TestCreatePersonNoForename(t *testing.T) {
	expectedError := "you must specify the Forename"
	personform := CreatePersonForm(expectedID, "", expectedSurname)
	if personform.Validate() {
		t.Errorf("Expected the validation to fail - no forename")
	} else {
		if personform.ErrorForField("Forename") != expectedError {
			t.Errorf("Expected \"%s\", got \"%s\"", expectedError,
				personform.ErrorForField("Forename"))
		}
	}
	errors := personform.FieldErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
}

// Create a personform containing a person with no surname, and validate it.
func TestCreatePersonNoSurname(t *testing.T) {
	expectedError := "you must specify the Surname"
	personform := CreatePersonForm(expectedID, expectedForename, "")
	if personform.Validate() {
		t.Errorf("Expected the validation to fail - no surname")
	} else {
		if personform.ErrorForField("Surname") != expectedError {
			t.Errorf("Expected \"%s\", got \"%s\"", expectedError,
				personform.ErrorForField("Surname"))
		}
	}
	errors := personform.FieldErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
}

func CreatePersonForm(id uint64, forename, surname string) ConcretePersonForm {
	var person model.GorpMysqlPerson
	person.SetID(id)
	person.SetForename(forename)
	person.SetSurname(surname)
	var personform ConcretePersonForm
	personform.SetPerson(&person)
	return personform
}
