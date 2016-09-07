package dbsession

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	personModel "github.com/goblimey/films/models/person"
	gorpModel "github.com/goblimey/films/models/person/gorpmysql"
	gorp "gopkg.in/gorp.v1"
	// This import must be present to satisfy a dependency in the GORP library.
	_ "github.com/go-sql-driver/mysql"
)

// The GorpMysqlDBSession type represents a MySQL database session accessed via GORP.
// It satisfies the DBSession interface.
type GorpMysqlDBSession struct {
	dbmap *gorp.DbMap
}

// MakeGorpMysqlDBSession is a factory function that creates a GorpMysqlDBSession and returns it as a pointer to a DBSession.
func MakeGorpMysqlDBSession() (DBSession, error) {
	log.SetPrefix("DBSessionFactory.MakeGorpMysqlDBSession() ")
	db, err := sql.Open("mysql", "webuser:secret@tcp(localhost:3306)/films")
	if err != nil {
		log.Printf("failed to get DB handle - %s\n" + err.Error())
		return nil, errors.New("failed to get DB handle - " + err.Error())
	}
	// check that the handle works
	err = db.Ping()
	if err != nil {
		log.Printf("cannot connect to DB.  %s\n", err.Error())
		return nil, err
	}
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	table := dbmap.AddTableWithName(gorpModel.GorpMysqlPerson{}, "people").SetKeys(true, "IDField")
	if table == nil {
		em := "cannot add table people"
		log.Println(em)
		return nil, errors.New(em)
	}

	table.ColMap("IDField").Rename("id")
	table.ColMap("ForenameField").Rename("forename")
	table.ColMap("SurnameField").Rename("surname")

	// Create any missing tables.
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		em := fmt.Sprintf("cannot create table - %s\n", err.Error())
		log.Printf("em")
		return nil, errors.New(em)
	}

	// Create a concrete DBSession and an interface reference to it.
	var session DBSession = &GorpMysqlDBSession{dbmap}

	// Return the interface reference.
	return session, nil
}

// StartTransaction starts a transaction.
func (dbs GorpMysqlDBSession) StartTransaction() (*gorp.Transaction, error) {
	return dbs.dbmap.Begin()
}

// Close closes the GORP DBMap and releases the database connection.
// Anything that opens a connection should call this method to close it.
func (dbs GorpMysqlDBSession) Close() {
	dbs.dbmap.Db.Close()
}

// FindAllPeople returns a slice of all valid Person records from the database in a
// (possibly empty) slice.  If the database lookup fails, the error is returned
// instead.
func (dbs GorpMysqlDBSession) FindAllPeople() ([]personModel.Person, error) {
	/*
	 * Get all Person records from the database into a slice, create a slice
	 * of the same size and copy the valid records into it.  Return the
	 * slice of valid records, which may be empty.  If the select fails, return
	 * the error.
	 */
	var GorpMysqlPersons []gorpModel.GorpMysqlPerson
	_, err := dbs.dbmap.Select(&GorpMysqlPersons, "select id, surname, forename from people")
	if err != nil {
		return nil, err
	}

	validPeople := make([]personModel.Person, len(GorpMysqlPersons))

	// Validate and copy the Person records
	next := 0 // Index of next validPeople entry
	for _, p := range GorpMysqlPersons {
		p.SetForename(strings.TrimSpace(p.Forename()))
		p.SetSurname(strings.TrimSpace(p.Surname()))
		if len(p.Forename()) > 0 && len(p.Surname()) > 0 {
			// This doesn't work - all entries end up containing the last added value
			// validPeople[next] = &person
			// We must clone the data instead
			person := personModel.Clone(&p)
			validPeople[next] = person
			next++
		}
	}

	return validPeople, nil
}

// FindPersonByID fetches the row from the people table with the given uint64 id. The
// data fetched may or may not be valid.  The method returns a Person containing
// that data, or an error message.
func (dbs GorpMysqlDBSession) FindPersonByID(id uint64) (personModel.Person, error) {
	m := "FindPersonByID()"
	log.Printf("%s: ID %d", m, id)
	var GorpMysqlPerson gorpModel.GorpMysqlPerson
	err := dbs.dbmap.SelectOne(&GorpMysqlPerson, "select id, surname, forename from people where id = ?", id)
	if err != nil {
		log.Printf("%s: %s", m, err.Error())
		return nil, err
	}
	log.Printf("%s: found person %s", m, GorpMysqlPerson.String())
	return &GorpMysqlPerson, nil
}
