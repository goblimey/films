package people

import (
	personModel "github.com/goblimey/films/models/person"
	"github.com/goblimey/films/utilities/dbsession"
)

// DAO is the interface defining Data Access Objects for the people table.
type DAO interface {

	SetSession(session dbsession.DBSession)
	
	/*
	FindAll() returns a pointer to a map of valid People indexed by ID.  Any
	invalid records are left out of the map
	*/
	FindAll() ([]personModel.Person, error)

	/*
	FindByid fetches the row from the people table with the given uint64 id. It validates that data
	and, if it's valid, uses it to create a Person and returns a pointer to it.  If the data is not
	valid the function returns an error message.
	*/
	FindByID(id uint64) (personModel.Person, error)

	/*
	FindByid fetches the row from the people table with the given string id. It validates that data
	and, if it's valid, uses it to create a Person and returns a pointer to it.  If the data is not
	valid the function returns an error message.  The ID in the database is numeric and the method
	checks that the given ID is also numeric before
	it makes the call.  If not, it returns an error.
	*/
	FindByIDStr(idStr string) (personModel.Person, error)

	/*
	Create takes a person and creates a record in the people table containing the same
	data and with an auto-incremented ID.  It returns a pointer to the resulting person
	or any error that the DB call supplies to it.
	*/
	Create(person personModel.Person) (personModel.Person, error)

	/*
	Update takes a person structure, updates the record in the people table with the
	same ID and returns the row count and error that the DB call supplies to it.  On
	a successful update, the number of rows returned should be 1.
	*/
	Update(person personModel.Person) (uint64, error)
	/*
	 * DeleteById takes the given uint64 ID and deletes the record with that ID from the people
	 * table.  The method returns the row count and error that the database supplies to it.  On
	 * a successful delete, it should return 1, having deleted one row.
	*/
	DeleteByID(id uint64) (int64, error)

	/*
	 * DeleteByIdStr takes the given String ID and deletes the record with that ID from the people
	 * table. The ID in the database is numeric and the method checks that the given ID is also
	 * numeric before it makes the call.  If not, it returns an error.  If the ID looks sensible,
	 * the function attempts the delete and returns the row count and error that the database
	 * supplies to it.  On a successful delete, it should return 1, having deleted one row.
	*/
	DeleteByIDStr(idStr string) (int64, error)
}