package manual

import (
	"errors"

	personModel "github.com/goblimey/films/models/person"
	"github.com/goblimey/films/utilities/dbsession"
)

// MockRepo is a hand-crafted mock which satisfies the Repository interface.
type MockRepo struct {
	PersonList []personModel.Person
}

var findAllCalled = false

// TestComplete checks that the methods have been called as expected
func (mr MockRepo) TestComplete() error {
	em := ""
	if !findAllCalled {
		em += "expected MockRepo.FindAll() to be called "
	}

	if len(em) > 0 {
		return errors.New(em)
	}

	// expectations satisfied
	return nil
}

// SetSession sets the database session
func (mr MockRepo) SetSession(session dbsession.DBSession) {
}

// FindAll returns a list of valid People - a list of Person objects each
// of which is valid according to Person.Validate()
func (mr MockRepo) FindAll() ([]personModel.Person, error) {
	findAllCalled = true
	if mr.PersonList == nil {
		// Return an error to test error handling
		return nil, errors.New("Test Error Message")
	}

	return mr.PersonList, nil

}

// FindByID fetches the row from the people table with the given uint64 id. It validates that data
// and, if it's valid, uses it to create a Person and returns a pointer to it.  If the data is not
// valid the function returns an error message.
func (mr MockRepo) FindByID(id uint64) (personModel.Person, error) {
	return nil, errors.New("FindById(): not expected this method to be called")
}

// FindByIDStr fetches the row from the people table with the given string id. It validates that data
// and, if it's valid, uses it to create a Person and returns a pointer to it.  If the data is not
// valid the function returns an error message.  The ID in the database is numeric and the method
// checks that the given ID is also numeric before
// it makes the call.  If not, it returns an error.
func (mr MockRepo) FindByIDStr(idStr string) (personModel.Person, error) {
	return nil, errors.New("FindByIDStr(): not expected this method to be called")

}

// Create takes a person, creates a record in the people table containing the same data
// and returns any error that the DB call supplies to it.  On a successful create, the
// error will be nil.
func (mr MockRepo) Create(person personModel.Person) (personModel.Person, error) {
	return nil, errors.New("Create(): not expected this method to be called")
}

// Update takes a person structure, updates the record in the people table
// with the same ID and returns the row count and error that the DB call supplies to it.
// On a successful update, the number of rows returned should be 1.
func (mr MockRepo) Update(person personModel.Person) (uint64, error) {
	return 0, errors.New("Update(): not expected this method to be called")
}

// DeleteByID takes the given uint64 ID and deletes the record with that ID from the people
// table.  The method returns the row count and error that the database supplies to it.  On
// a successful delete, it should return 1, having deleted one row.
func (mr MockRepo) DeleteByID(id uint64) (int64, error) {
	return 0, errors.New("DeleteByID(): not expected this method to be called")
}

// DeleteByIDStr takes the given String ID and deletes the record with that ID from the people
// table. The ID in the database is numeric and the method checks that the given ID is also
// numeric before it makes the call.  If not, it returns an error.  If the ID looks sensible,
// the function attempts the delete and returns the row count and error that the database
// supplies to it.  On a successful delete, it should return 1, having deleted one row.
func (mr MockRepo) DeleteByIDStr(idStr string) (int64, error) {
	return 0, errors.New("DeleteByIDStr(): not expected this method to be called")
}
