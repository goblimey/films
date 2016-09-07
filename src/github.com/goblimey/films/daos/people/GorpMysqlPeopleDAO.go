// Package people provides Create, Read, Update and Delete (CRUD) operations on the
// people resource.  That resource is referenced via a database session that is
// supplied by the parent.  For example it could be a MySQL table accessed via GORP,
// but it could also be a mock session.
//
// The GorpMysqlDAO satisfies the DAO interface.
package people

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	personModel "github.com/goblimey/films/models/person"
	gorpPersonModel "github.com/goblimey/films/models/person/gorpmysql"
	"github.com/goblimey/films/utilities/dbsession"
)

// GorpMysqlDAO satifies the DAO interface.
type GorpMysqlDAO struct {
	session dbsession.DBSession
}

// MakeDAO is a factory function that creates a GorpMysqlDAO and returns it as a DAO.
func MakeDAO(session dbsession.DBSession) DAO {
	var DAO DAO = &GorpMysqlDAO{session}
	return DAO
}

// SetSession sets the session.
func (gmpd *GorpMysqlDAO) SetSession(session dbsession.DBSession) {
	gmpd.session = session
}

// FindAll returns a list of all valid Person records from the database in a slice.
// The result may be an empty slice.  If the database lookup fails, the error is
// returned instead.
func (gmpd GorpMysqlDAO) FindAll() ([]personModel.Person, error) {
	m := "FindAll()"
	log.Printf("%s:\n", m)
	people, err := gmpd.session.FindAllPeople()
	return people, err
}

// FindByID fetches the row from the people table with the given uint64 id. It
// validates that data and, if it's valid, returns the person.  If the data is not
// valid the function returns an error message.
func (gmpd GorpMysqlDAO) FindByID(id uint64) (personModel.Person, error) {
	m := "FindByID()"
	log.Printf("%s: ID %d", m, id)

	var person personModel.Person
	person, err := gmpd.session.FindPersonByID(id)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(person.Forename())) < 1 {
		return nil, errors.New("invalid person - no forename")
	}
	if len(strings.TrimSpace(person.Surname())) < 1 {
		return nil, errors.New("invalid person - no surname")
	}
	return person, nil
}

// FindByIDStr fetches the row from the people table with the given string id. It
// validates that data and, if it's valid, returns the person.  If the data is not valid
// the function returns an errormessage.  The ID in the database is numeric and the method
// checks that the given ID is also numeric before it makes the call.  This avoids hitting
// the DB when the id is obviously junk.
func (gmpd GorpMysqlDAO) FindByIDStr(idStr string) (personModel.Person, error) {
	m := "FindByIDStr()"
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		em := fmt.Sprintf("ID %s is not an unsigned integer", idStr)
		log.Printf("%s: %s", m, em)
		return nil, fmt.Errorf("ID %s is not an unsigned integer", idStr)
	}
	return gmpd.FindByID(id)
}

// Create takes a person, creates a record in the people table containing the same
// data with an auto-incremented ID and returns any error that the DB call returns.
// On a successful create, the method returns the created person, including
// the assigned ID.  This is all done within a transaction to ensure atomicity.
func (gmpd GorpMysqlDAO) Create(person personModel.Person) (personModel.Person, error) {
	m := "Create()"
	log.Printf("%s:", m)
	tx, err := gmpd.session.StartTransaction()
	if err != nil {
		log.Printf("%s: %s", m, err.Error())
		return nil, err
	}
	person.SetID(0) // provokes the auto-increment
	err = tx.Insert(person)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Printf("%s: created person %s", m, person.String())
	return person, nil
}

// Update takes a person record, updates the record in the people table with the same ID
// and returns the updated person or any error that the DB call supplies to it.  The update
// is done within a transaction
func (gmpd GorpMysqlDAO) Update(person personModel.Person) (uint64, error) {
	m := "Update()"
	tx, err := gmpd.session.StartTransaction()
	if err != nil {
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}
	rowsUpdated, err := tx.Update(person)
	if err != nil {
		tx.Rollback()
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}
	if rowsUpdated != 1 {
		tx.Rollback()
		em := fmt.Sprintf("update failed - %d rows would have been updated, expected 1", rowsUpdated)
		log.Printf("%s: %s", m, em)
		return 0, errors.New(em)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}

	// Success!
	return 1, nil
}

// DeleteByID takes the given uint64 ID and deletes the record with that ID from the people table.
// The function returns the row count and error that the database supplies to it.  On a successful
// delete, it should return 1, having deleted one row.
func (gmpd GorpMysqlDAO) DeleteByID(id uint64) (int64, error) {
	m := "DeleteByID()"
	log.Printf("%s: ID %d", m, id)
	// Need a Person record for the delete method, so fake one up.
	var person gorpPersonModel.GorpMysqlPerson
	person.SetID(id)
	tx, err := gmpd.session.StartTransaction()
	if err != nil {
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}
	rowsDeleted, err := tx.Delete(&person)
	if err != nil {
		tx.Rollback()
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}
	if rowsDeleted != 1 {
		tx.Rollback()
		em := fmt.Sprintf("delete failed - %d rows would have been deleted, expected 1", rowsDeleted)
		log.Printf("%s: %s", m, em)
		return 0, errors.New(em)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("%s: %s", m, err.Error())
		return 0, err
	}
	if err != nil {
		log.Printf("%s: %s", m, err.Error())
	}
	return rowsDeleted, nil
}

// DeleteByIDStr takes the given String ID and deletes the record with that ID from the people table.
// The ID in the database is numeric and the method checks that the given ID is also numeric before
// it makes the call.  If not, it returns an error.  If the ID looks sensible, the function attempts
// the delete and returns the row count and error that the database supplies to it.  On a successful
// delete, it should return 1, having deleted one row.
func (gmpd GorpMysqlDAO) DeleteByIDStr(idStr string) (int64, error) {
	m := "DeleteByIDStr()"
	log.Printf("%s: ID %s", m, idStr)
	// Check the id.
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		em := fmt.Sprintf("ID %s is not an unsigned integer", idStr)
		log.Printf("%s: %s", m, em)
		return 0, errors.New(em)
	}
	return gmpd.DeleteByID(id)
}
