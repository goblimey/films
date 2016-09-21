package gorpmysql

import (
	"testing"

	personModel "github.com/goblimey/films/models/person"
)

var expectedID uint64 = 2
var expectedForename = "foo"
var expectedSurname = "bar"

var person personModel.Person

func init() {
	person = MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
}

func TestIntegrationCreateGorpMysqlPersonCheckID(t *testing.T) {
	if person.ID() != expectedID {
		t.Errorf("expected ID to be %d actually %d", expectedID, person.ID())
	}
}

func TestIntegrationCreatePersonCheckForename(t *testing.T) {
	if person.Forename() != expectedForename {
		t.Errorf("expected forename to be %s actually %s", expectedForename, person.Forename())
	}
}

func TestIntegrationCreatePersonCheckSurname(t *testing.T) {
	if person.Surname() != expectedSurname {
		t.Errorf("expected surname to be %s actually %s", expectedSurname, person.Surname())
	}
}
