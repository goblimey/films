package people

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	personModel "github.com/goblimey/films/models/person/gorpmysql"
	dbsession "github.com/goblimey/films/utilities/dbsession"
)

// This is an integration test for the GorpMysqlDAO connecting to a MySQL DB via GORP.

var expectedForename1 = "foo"
var expectedSurname1 = "bar"
var expectedForename2 = "foobar"
var expectedSurname2 = "barfoo"

// Create a person in the database, read it back, test the contents.
func TestIntCreatePersonStoreFetchBackAndCheckContents(t *testing.T) {
	log.SetPrefix("TestIntegrationCreatePersonAndCheckContents")
	// Create a dao containing a session
	dbsession, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer dbsession.Close()

	dao := MakeDAO(dbsession)

	clearDown(dao, t)

	p := personModel.MakeInitialisedPerson(0, expectedForename1, expectedSurname1)

	// Store the person in the DB
	person, err := dao.Create(p)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("created person %s\n", person.String())

	retrievedPerson, err := dao.FindByID(person.ID())
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("retrieved person %s\n", retrievedPerson.String())

	if retrievedPerson.ID() != person.ID() {
		t.Errorf("expected ID to be %d actually %d", person.ID(),
			retrievedPerson.ID())
	}
	if retrievedPerson.Forename() != expectedForename1 {
		t.Errorf("expected forename to be %s actually %s", expectedForename1,
			retrievedPerson.Forename())
	}
	if retrievedPerson.Surname() != expectedSurname1 {
		t.Errorf("expected surname to be %s actually %s", expectedSurname1,
			retrievedPerson.Surname())
	}

	// Delete person and check response
	id := person.ID()
	rows, err := dao.DeleteByID(person.ID())
	if err != nil {
		t.Errorf(err.Error())
	}
	if rows != 1 {
		t.Errorf("expected delete to return 1, actual %d", rows)
	}
	log.Printf("deleted person with ID %d\n", id)
	clearDown(dao, t)
}

// Create a two person records in the DB, read them back and check the fields
func TestIntCreateTwoPersonsAndReadBack(t *testing.T) {
	log.SetPrefix("TestCreatePersonAndReadBack")
	// Create a dao containing a session
	dbsession, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer dbsession.Close()

	dao := MakeDAO(dbsession)

	clearDown(dao, t)

	//Create two people
	p1 := personModel.MakeInitialisedPerson(0, expectedForename1, expectedSurname1)
	person1, err := dao.Create(p1)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Printf("person1 %s", person1.String())
	p2 := personModel.MakeInitialisedPerson(0, expectedForename2, expectedSurname2)
	person2, err := dao.Create(p2)
	if err != nil {
		t.Errorf(err.Error())
	}

	// read all the people in the DB - expect just the two we created
	people, err := dao.FindAll()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(people) != 2 {
		t.Errorf("expected 2 rows, actual %d", len(people))
	}

	matches := 0
	for _, person := range people {
		switch person.Forename() {
		case expectedForename1:
			if person.Surname() == expectedSurname1 {
				matches++
			} else {
				t.Errorf("expected surname to be %s actually %s", expectedSurname1, person.Surname())
			}
		case expectedForename2:
			if person.Surname() == expectedSurname2 {
				matches++
			} else {
				t.Errorf("expected forename to be %s actually %s", expectedForename2, person.Forename())
			}
		default:
			t.Errorf("unexpected forename - %s", person.Forename())
		}
	}

	// We should have just the records we created
	if matches != 2 {
		t.Errorf("expected two matches, actual %d", matches)
	}

	// Find each of the records by ID and check the fields
	log.Printf("finding person %d", person1.ID())
	person1Returned, err := dao.FindByID(person1.ID())
	if err != nil {
		t.Errorf(err.Error())
	}
	if person1Returned.Forename() != expectedForename1 {
		t.Errorf("expected forename to be %s actually %s",
			expectedForename1, person1Returned.Forename())
	}
	if person1Returned.Surname() != expectedSurname1 {
		t.Errorf("expected surname to be %s actually %s",
			expectedSurname1, person1Returned.Surname())
	}

	var IDStr = strconv.FormatUint(person2.ID(), 10)
	person2Returned, err := dao.FindByIDStr(IDStr)
	if err != nil {
		t.Errorf(err.Error())
	}

	log.Printf("found person %s", person2Returned.String())

	if person2Returned.Forename() != expectedForename2 {
		t.Errorf("expected forename to be %s actually %s",
			expectedForename2, person2Returned.Forename())
	}
	if person2Returned.Surname() != expectedSurname2 {
		t.Errorf("expected surname to be %s actually %s",
			expectedSurname2, person2Returned.Surname())
	}

	clearDown(dao, t)
}

// Create two people, remove one, check that we get back Just the other
func TestIntCreateTwoPeopleAndDeleteOneByIDStr(t *testing.T) {
	log.SetPrefix("TestIntegrationCreateTwoPeopleAndDeleteOneByIDStr")
	// Create a dao containing a session
	dbsession, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer dbsession.Close()

	dao := MakeDAO(dbsession)

	clearDown(dao, t)

	// Create two people
	p1 := personModel.MakeInitialisedPerson(0, expectedForename1, expectedSurname1)
	person1, err := dao.Create(p1)
	if err != nil {
		t.Errorf(err.Error())
	}

	p2 := personModel.MakeInitialisedPerson(0, expectedForename2, expectedSurname2)
	person2, err := dao.Create(p2)
	if err != nil {
		t.Errorf(err.Error())
	}

	var IDStr = fmt.Sprintf("%d", person1.ID())
	rows, err := dao.DeleteByIDStr(IDStr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if rows != 1 {
		t.Errorf("expected one record to be deleted, actually %d", rows)
	}

	// We should have one record in the DB and it should match person2
	people, err := dao.FindAll()
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(people) != 1 {
		t.Errorf("expected one record, actual %d", len(people))
	}

	for _, person := range people {
		if person.ID() != person2.ID() {
			t.Errorf("expected id to be %d actually %d",
				person2.ID(), person.ID())
		}
		if person.Forename() != expectedForename2 {
			t.Errorf("expected forename to be %s actually %s",
				expectedForename2, person.Forename())
		}
		if person.Surname() != expectedSurname2 {
			t.Errorf("TestCreateTwoPeopleAndDeleteOneByIDStr(): expected surname to be %s actually %s",
				expectedSurname2, person.Surname())
		}
	}

	clearDown(dao, t)
}

// Create a person record, update the record, read it back and check that it's updated
func TestIntCreatePersonAndUpdate(t *testing.T) {
	log.SetPrefix("TestIntegrationCreatePersonAndUpdate")
	// Create a dao containing a session
	dbsession, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer dbsession.Close()

	dao := MakeDAO(dbsession)

	clearDown(dao, t)

	// Create a person in the DB.
	p := personModel.MakeInitialisedPerson(0, expectedSurname1, expectedForename1)
	person, err := dao.Create(p)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Update the person in the DB.
	person.SetForename(expectedForename2)
	person.SetSurname(expectedSurname2)
	rows, err := dao.Update(person)
	if err != nil {
		t.Errorf(err.Error())
	}
	if rows != 1 {
		t.Errorf("expected 1 row to be updated, actually %d rows", rows)
	}

	// fetch the updated record back and check it.
	personfetched, err := dao.FindByID(person.ID())
	if err != nil {
		t.Errorf(err.Error())
	}

	if personfetched.Forename() != expectedForename2 {
		t.Errorf("expected forename to be %s actually %s",
			expectedForename2, personfetched.Forename())
	}
	if person.Surname() != expectedSurname2 {
		t.Errorf("expected surname to be %s actually %s",
			expectedSurname2, personfetched.Surname())
	}

	clearDown(dao, t)
}

// clearDown() - helper function to remove all people from the DB
func clearDown(dao DAO, t *testing.T) {
	people, err := dao.FindAll()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	for _, person := range people {
		rows, err := dao.DeleteByID(person.ID())
		if err != nil {
			t.Errorf(err.Error())
			continue
		}
		if rows != 1 {
			t.Errorf("while clearing down, expected 1 row, actual %d", rows)
		}
	}
}
