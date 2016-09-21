package dbsession

import (
	gorp "gopkg.in/gorp.v1"
	personModel "github.com/goblimey/films/models/person"
)

// DBSession represents a database session.
type DBSession interface {

	/*
	Start a new transaction.  A transaction is a resource overhead and the caller should
	call Close() when it's finished to release this resource.
	*/
	StartTransaction() (*gorp.Transaction, error)
	
	// Close the DBSession and release the resources associated with it.
	Close()

	/*
	FindAllPeople() gets all records in the people table (whether valid or not) and returns a 
	pointer to a slice containing them.  The method does not create an explicit transaction.
	*/
	FindAllPeople() ([]personModel.Person, error)

	/*
	 FindPersonByid fetches the row from the people table with the given uint64 id. The
	 data fetched may or may not be valid.  The method returns a Person containing
	 that data, or an error message.
	*/
	FindPersonByID(id uint64) (personModel.Person, error)
}