/*
 * Package people provides the controller for the people resource.  It provides a set of action
 * functions that are triggered by HTTP requests and implement the Create, Read, Update and
 * Delete (CRUD) operations on the people resource:

 *    GET people/ - runs Index() to list all people
 *    GET people/n - runs Show() to display the details of the person with ID n
 *    GET people/create - runs New() to display the page to create a person using any data in the form to pre-populate it
 *    PUT people/n - runs Create() to create a new person using the data in the supplied form
 *    GET people/n/edit - runs Edit() to display the page to edit the person with ID n, using any data in the form to pre-populate it
 *    PUT people/n - runs Update() to update the person with ID n using the data in the form
 *    DELETE people/n - runs Delete() to delete te person with id n
 *
 * The requests follow the REST model and therefore avoid carrying data such as the ID in HTTP
 * parameters.
 *
 * The REST model uses HTTP requests such as PUT and DELETE.  The standard browsers do not support
 * these operrations, so they are implemented using a POST request with a parameter "_method"
 * defining the operation.  (A post with a parameter "_method=PUT" simulates a PUT, and so on.)
 */

package people

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	daos "github.com/goblimey/films/daos"
	pd "github.com/goblimey/films/daos/people"
	forms "github.com/goblimey/films/forms/people"
	personModel "github.com/goblimey/films/models/person/gorpmysql"
	request "github.com/goblimey/films/retrofit/restful/request"
	response "github.com/goblimey/films/retrofit/restful/response"
	"github.com/goblimey/films/template"
	"github.com/goblimey/films/utilities"
	"github.com/goblimey/films/utilities/dbsession"
)

/*
 * The deleteRequestRE is the regular expression for the URI of a delete
 * request containing a numeric ID - for example: "/people/1/delete".
 */
var deleteRequestRE = regexp.MustCompile(`^/people/[0-9]+/delete$`)

/*
 * The showRequestRE is the regular expression for the URI of a show
 * request containing a numeric ID - for example: "/people/1".
 */
var showRequestRE = regexp.MustCompile(`^/people/[0-9]+$`)

/*
 * The editRequestRE is the regular expression for the URI of an edit
 * request containing a numeric ID - for example: "/people/1/edit".
 */
var editRequestRE = regexp.MustCompile(`^/people/[0-9]+/edit$`)

/*
 * The updateRequestRE is the regular expression for the URI of an update
 * request containing a numeric ID - for example: "/people/1".  The URI
 * is the same as for the show request, but we give it a different name for
 * clarity.
 */
var updateRequestRE = showRequestRE

// Controller is the MVC controller class for the people resource.
type Controller struct{}

// Marshall calls the appropriate request handler for the request
func (c Controller) Marshall(req request.Request, resp response.Response,
	templates map[string]template.Template) {

	log.SetPrefix("people.Controller.Marshall() ")

	uri := req.HTTPRequest().URL().RequestURI()
	log.Printf("uri %s", uri)

	method := req.HTTPRequest().Method()
	if method == "POST" {
		// handle simulated PUT via _method parameter
		simMethod := req.HTTPRequest().FormValue("_method")
		if simMethod == "PUT" || simMethod == "DELETE" {
			method = simMethod
		}
	}
	log.Printf("method %s uri %s", method, uri)

	// Create a dao containing a session (GORP implementation).
	session, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		em := fmt.Sprintf("cannot create a session - %s", err.Error())
		log.Println(em)
		utilities.Dead(resp)
		return
	}
	defer session.Close()

	// Create a DAO service
	pdi := new(pd.GorpMysqlDAO)
	pdi.SetSession(session)
	var dao pd.DAO = pdi
	daoService := daos.MakeDAOService(dao)

	// Call the appropriate handler for the request

	switch method {

	case "GET":

		if uri == "/people" {
			// "GET http://server:port/people" - fetch all the valid people
			// records and display them.
			var form forms.ConcreteListForm
			Index(req, resp, daoService, &form, templates)

		} else if editRequestRE.MatchString(uri) {

			// "GET http://server:port/people/1/edit" - fetch the people record
			// given by the ID in the request and display the form to edit it.
			var form forms.ConcretePersonForm
			Edit(req, resp, daoService, &form, templates)

		} else if uri == "/people/create" {

			// "GET http://server:port/people/create" - display the form to
			// create a new people record.
			var form forms.ConcretePersonForm
			// Create an empty person to get started.
			person := personModel.MakePerson()
			form.SetPerson(person)
			New(req, resp, daoService, &form, templates)

		} else if showRequestRE.MatchString(uri) {

			// "GET http://server:port/people/1" - fetch the people record with ID 1 and display it.
			var form forms.ConcretePersonForm
			Show(req, resp, daoService, &form, templates)
		}

	case "PUT":
		if updateRequestRE.MatchString(uri) {

			// POST http://server:port/people/1" - update the people record with
			// the given ID from the URI using the form data in the body.
			form := getPersonFormFromRequest(req, resp, daoService, templates)
			Update(req, resp, daoService, form, templates)

		} else if uri == "/people" {

			// POST http://server:port/people" - create a new people record from
			// the form data in the body.
			form := getPersonFormFromRequest(req, resp, daoService, templates)
			Create(req, resp, daoService, form, templates)
		}

	case "DELETE":
		if deleteRequestRE.MatchString(uri) {

			// "POST http://server:port/people/1/delete" - delete the people
			// record with the ID given in the request.
			Delete(req, resp, daoService, templates)
		}

	default:
		em := fmt.Sprintf("unexpected HTTP method %v", method)
		log.Println(em)
		errorHandler(req, resp, daoService, em, templates)
	}
}

// Index fetches a list of all valid people and displays the index page.
func Index(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.ListForm,
	templates map[string]template.Template) {

	log.SetPrefix("Index()")

	listPeople(req, resp, daoService, form, templates)
	return
}

// Show displays the details of the person with the ID given in the URI.
func Show(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.PersonForm,
	templates map[string]template.Template) {

	log.SetPrefix("Show()")

	id := req.PathParameter("id")
	log.Printf("id=%s", id)

	// Get the details of the person with the given ID.
	dao := daoService.DAO()
	person, err := dao.FindByIDStr(id)
	if err != nil {
		// no such person.  Display index page with error message
		em := "no such person"
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}

	form.SetPerson(person)

	page := templates["Show"]
	if page == nil {
		em := fmt.Sprintf("internal error displaying Show page - no HTML template")
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}

	err = page.Execute(resp.Response().ResponseWriter, form)
	if err != nil {
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	return
}

// New displays the page to create a new person,
func New(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.PersonForm,
	templates map[string]template.Template) {

	log.SetPrefix("New()")

	// Display the page.
	page := templates["Create"]
	if page == nil {
		em := fmt.Sprintf("internal error displaying Create page - no HTML template")
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	err := page.Execute(resp.Response().ResponseWriter, form)
	if err != nil {
		log.Printf("error displaying new page - %s", err.Error())
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		errorHandler(req, resp, daoService, em, templates)
		return
	}
}

// Create creates a new person using the data from the HTTP form displayed
// by a previous NEW request.
func Create(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.PersonForm,
	templates map[string]template.Template) {

	log.SetPrefix("Create()")

	if !(form.Validate()) {
		// validation errors.  Return to create screen with error messages in the form data
		page := templates["Create"]
		if page == nil {
			em := fmt.Sprintf("internal error displaying Create page - no HTML template")
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
			return
		}
		err := page.Execute(resp.Response().ResponseWriter, &form)
		if err != nil {
			em := fmt.Sprintf("Internal error while preparing create form after failed validation - %s",
				err.Error())
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
			return
		}
		return
	}

	// Create a person in the database using the validated data in the form
	dao := daoService.DAO()
	createdPerson, err := dao.Create(form.Person())
	if err != nil {
		// Failed to create person.  Display index page with error message.
		em := fmt.Sprintf("Could not create person %s - %s", form.Person().String(), err.Error())
		errorHandler(req, resp, daoService, em, templates)
		return
	}

	// Success! Person created.  Display index page with confirmation notice
	notice := fmt.Sprintf("created new person %s", createdPerson.String())
	log.Printf("%s\n", notice)
	var cpf forms.ConcreteListForm
	var listForm forms.ListForm = &cpf
	listForm.SetNotice(notice)
	listPeople(req, resp, daoService, listForm, templates)
	return
}

// Edit fetches the data for the people record with the given ID and displays
// the edit page, populated with that data.
func Edit(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.PersonForm,
	templates map[string]template.Template) {

	log.SetPrefix("Edit() ")

	log.Println("parsing form")
	err := req.HTTPRequest().ParseForm()
	if err != nil {
		// failed to parse form
		em := fmt.Sprintf("cannot parse form - %s", err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	// Get the ID of the person
	id := req.PathParameter("id")

	dao := daoService.DAO()
	// Get the existing data for the person
	person, err := dao.FindByIDStr(id)
	if err != nil {
		// No such person.  Display index page with error message.
		em := err.Error()
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	/*
	 * Got the person with the given ID.  Put it into the form and validate it.
	 * If the data is invalid, continue - the user may be trying to fix it.
	 */
	form.SetPerson(person)
	if !form.Validate() {
		em := fmt.Sprintf("invalid record in the people database - %s",
			person.String())
		log.Printf("%s\n", em)
	}

	// Display the edit page
	page := templates["Edit"]
	if page == nil {
		em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	log.Printf("Execute\n")
	err = page.Execute(resp.Response().ResponseWriter, form)
	if err != nil {
		// error while preparing edit page
		log.Printf("%s: error displaying edit page - %s", err.Error())
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		errorHandler(req, resp, daoService, em, templates)
	}
}

// Update responds to a PUT request.  For example:
// PUT /people/1
// It's invoked by the form displayed by a previous Edit request.  If the ID in the URI is
// valid and the request parameters from the form specify valid people data, it updates the
// record and displays the index page with a confirmation message, otherwise it displays
// the edit page again with the given data and some error messages.
func Update(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.PersonForm,
	templates map[string]template.Template) {

	log.SetPrefix("Update() ")

	// Get the person specified in the form from the DB.
	// (which also validates the id in the form).
	log.Printf("form=%v\n", form)
	if form.Person() == nil {
		em := fmt.Sprint("internal error - form should contain an updated person record")
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	dao := daoService.DAO()
	person, err := dao.FindByID(form.Person().ID())
	if err != nil {
		// There is no person with this ID.  The ID is chosen by the user from a
		// supplied list and it should always be valid, so there's something screwy
		// going on.  Display the index page with an error message.
		em := fmt.Sprintf("error searching for person with id %s - %s",
			form.Person().ID(), err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}

	// We have a matching person from the DB.
	log.Printf("got person %v\n", person)

	// Validate the new version of the person in the form.
	if !form.Validate() {
		// The data is invalid.  The validator has set error messages.  Return to
		// the edit screen.
		page := templates["Edit"]
		if page == nil {
			em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
			return
		}
		err = page.Execute(resp.Response().ResponseWriter, form)
		if err != nil {
			log.Printf("%s: error displaying edit page - %s", err.Error())
			em := fmt.Sprintf("error displaying page - %s", err.Error())
			errorHandler(req, resp, daoService, em, templates)
			return
		}
		return
	}

	// we have a valid record and valid new values.  Update.
	person.SetForename(form.Person().Forename())
	person.SetSurname(form.Person().Surname())
	log.Printf("updating person to %v\n", person)
	_, err = dao.Update(person)
	if err != nil {
		// The commit failed.  Display the edit page with an error message
		em := fmt.Sprintf("Could not update person - %s", err.Error())
		log.Printf("%s\n", em)
		form.SetErrorMessage(em)

		page := templates["Edit"]
		if page == nil {
			em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
			return
		}
		err = page.Execute(resp.Response().ResponseWriter, form)
		if err != nil {
			// Error while recovering from another error.  This is looking like a habit!
			em := fmt.Sprintf("Internal error while preparing edit page after failing to update person in DB - %s", err.Error())
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
		} else {
			return
		}
	}

	// Success!  Display the index page with a confirmation notice
	notice := fmt.Sprintf("updated person %s", form.Person().String())
	log.Printf("%s:\n", notice)
	var cpf forms.ConcreteListForm
	var listForm forms.ListForm = &cpf
	listForm.SetNotice(notice)
	listPeople(req, resp, daoService, listForm, templates)
	return
}

// Delete reponds to a DELETE request and deletes the record with the given ID,
// eg DELETE http://server:port/people/1.
func Delete(req request.Request, resp response.Response,
	daoService daos.DAOService, templates map[string]template.Template) {

	log.SetPrefix("Delete()")

	err := req.HTTPRequest().ParseForm()
	if err != nil {
		// failed - form does not parse
		em := fmt.Sprintf("Internal error - %s", err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	method := req.HTTPRequest().FormValue("_method")
	if "DELETE" != method {
		// failed - _method param is not DELETE
		em := fmt.Sprintf("Internal error - request type %s must be DELETE", method)
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	id := req.PathParameter("id")

	dao := daoService.DAO()
	// Attempt the delete
	_, err = dao.DeleteByIDStr(id)
	if err != nil {
		// failed - cannot delete person
		em := fmt.Sprintf("Cannot delete person with id %s - %s", id, err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return
	}
	// Success - person deleted.  Display the index view with a notification.
	var cpf forms.ConcreteListForm
	var listForm forms.ListForm = &cpf
	notice := fmt.Sprintf("deleted person with ID %s", id)
	log.Printf("%s:\n", notice)
	listForm.SetNotice(notice)
	listPeople(req, resp, daoService, listForm, templates)
	return
}

// getPersonFormFromRequest gets the person data from the request, creates a
// GorpMySQLPerson and returns it in a PersonForm.
func getPersonFormFromRequest(req request.Request, resp response.Response,
	daoService daos.DAOService, templates map[string]template.Template) forms.PersonForm {

	log.SetPrefix("getPersonFormFromRequest() ")

	err := req.HTTPRequest().ParseForm()
	if err != nil {
		em := fmt.Sprintf("cannot parse form - %s", err.Error())
		log.Printf("%s\n", em)
		errorHandler(req, resp, daoService, em, templates)
		return nil
	}
	var form forms.ConcretePersonForm
	var person personModel.GorpMysqlPerson
	idStr := req.PathParameter("id")
	if idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			em := fmt.Sprintf("invalid id %v in request - should be numeric", idStr)
			log.Printf("%s\n", em)
			errorHandler(req, resp, daoService, em, templates)
			return nil
		}
		person.SetID(id)
	}
	person.SetForename(strings.TrimSpace(req.HTTPRequest().FormValue("forename")))
	person.SetSurname(strings.TrimSpace(req.HTTPRequest().FormValue("surname")))
	form.SetPerson(&person)
	log.Printf("form %s\n", form.String())
	return &form
}

/*
 * The listPeople helper function fetches a list of people and displays the
 * index page.  It's used to fulfil an index request but the index page is
 * also used as the last page of a sequence of requests (for example new,
 * create, index).  If the sequence was successful, the form may contain a
 * confirmation note.  If the sequence failed, the form should contain an error
 * message.
 */
func listPeople(req request.Request, resp response.Response,
	daoService daos.DAOService, form forms.ListForm,
	templates map[string]template.Template) {

	log.SetPrefix("Controller.listPeople() ")

	dao := daoService.DAO()
	peopleList, err := dao.FindAll()
	if err != nil {
		em := fmt.Sprintf("error getting the list of people - %s", err.Error())
		log.Printf("%s\n", em)
		form.SetErrorMessage(em)
	} else {
		log.Printf("%d people", len(peopleList))
		if len(peopleList) <= 0 {
			form.SetNotice("there are no people currently set up")
		}
	}
	form.SetPeople(peopleList)

	// Display the index page
	page := templates["Index"]
	if page == nil {
		utilities.Dead(resp)
		return
	}
	err = page.Execute(resp.Response().ResponseWriter, form)
	if err != nil {
		/*
		 * Error while displaying the index page.  We handle most internal
		 * errors by displaying the index page.  That's just failed, so
		 * fall back to the static error page.
		 */
		log.Printf(err.Error())
		page = templates["Error"]
		if page == nil {
			utilities.Dead(resp)
			return
		}
		err = page.Execute(resp.Response().ResponseWriter, form)
		if err != nil {
			// Can't display the static error page either.  Bale out.
			em := fmt.Sprintf("fatal error - failed to display error page for error %s\n", err.Error())
			log.Printf(em)
			panic(em)
		}
		return
	}
}

// The errorHandler() helper function displays the index page with an error message
func errorHandler(req request.Request, resp response.Response,
	daoService daos.DAOService, errormessage string,
	templates map[string]template.Template) {

	var cpf forms.ConcreteListForm
	var listForm forms.ListForm = &cpf
	listForm.SetErrorMessage(errormessage)
	listPeople(req, resp, daoService, listForm, templates)
}
