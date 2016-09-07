package person

import (
	"testing"
)

func TestUnitCreateConcretePersonCheckID(t *testing.T) {
	var expectedID uint64 = 2
	var expectedForename string = "foo"
	var expectedSurname string = "bar"
	person := MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	if person.ID() != expectedID {
		t.Errorf("expected ID to be %d actually %d", expectedID, person.ID())
	}
}

func TestUnitCreatePersonCheckForename(t *testing.T) {
	var expectedID uint64 = 2
	var expectedForename string = "foo"
	var expectedSurname string = "bar"
	person := MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	if person.Forename() != expectedForename {
		t.Errorf("expected forename to be %s actually %s", expectedForename, person.Forename())
	}
}

func TestUnitCreatePersonCheckSurname(t *testing.T) {
	var expectedID uint64 = 2
	var expectedForename string = "foo"
	var expectedSurname string = "bar"
	person := MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	if person.Surname() != expectedSurname {
		t.Errorf("expected surname to be %s actually %s", expectedSurname, person.Surname())
	}
}
