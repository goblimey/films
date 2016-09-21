// Package people provides the controller for the people resource.  It provides a
// set of action functions that are triggered by HTTP requests and implement the
// Create, Read, Update and Delete (CRUD) operations on the people resource:
//
//    GET people/ - runs Index() to list all people
//    GET people/n - runs Show() to display the details of the person with ID n
//    GET people/create - runs New() to display the page to create a person using any data in the form to pre-populate it
//    PUT people/n - runs Create() to create a new person using the data in the supplied form
//    GET people/n/edit - runs Edit() to display the page to edit the person with ID n, using any data in the form to pre-populate it
//    PUT people/n - runs Update() to update the person with ID n using the data in the form
//    DELETE people/n - runs Delete() to delete the person with id n

package people

import (
	"fmt"
	"log"

	restful "github.com/emicklei/go-restful"
	forms "github.com/goblimey/films/forms/people"
	"github.com/goblimey/films/services"
	"github.com/goblimey/films/utilities"
)

type Controller struct {
	services services.Services
}

// MakeController is a factory that creates a people controller
func MakeController(services services.Services) Controller {
	var controller Controller
	controller.SetServices(services)
	return controller
}

// Index fetches a list of all valid people and displays the index page.
func (c Controller) Index(req *restful.Request, resp *restful.Response,
	form forms.ListForm) {

	log.SetPrefix("Index()")

	listPeople(req, resp, form, c.services)
	return
}

// Show displays the details of the person with the ID given in the URI.
func (c Controller) Show(req *restful.Request, resp *restful.Response,
	form forms.PersonForm) {

	log.SetPrefix("Show()")

	dao := c.services.GetPeopleRepository()

	// Get the details of the person with the given ID.
	person, err := dao.FindByID(form.Person().ID())
	if err != nil {
		// no such person.  Display index page with error message
		em := "no such person"
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}

	// The person in the form contains just an ID.  Replace it with the
	// complete person record that we just fetched.
	form.SetPerson(person)

	page := c.services.Template("Show")
	if page == nil {
		em := fmt.Sprintf("internal error displaying Show page - no HTML template")
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}

	err = page.Execute(resp.ResponseWriter, form)
	if err != nil {
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	return
}

// New displays the page to create a new person,
func (c Controller) New(req *restful.Request, resp *restful.Response,
	form forms.PersonForm) {

	log.SetPrefix("New()")

	// Display the page.
	page := c.services.Template("Create")
	if page == nil {
		em := fmt.Sprintf("internal error displaying Create page - no HTML template")
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	err := page.Execute(resp.ResponseWriter, form)
	if err != nil {
		log.Printf("error displaying new page - %s", err.Error())
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		c.ErrorHandler(req, resp, em)
		return
	}
}

// Create creates a new person using the data from the HTTP form displayed
// by a previous NEW request.
func (c Controller) Create(req *restful.Request, resp *restful.Response,
	form forms.PersonForm) {

	log.SetPrefix("Create()")

	if !(form.Validate()) {
		// validation errors.  Return to create screen with error messages in the form data
		page := c.services.Template("Create")
		if page == nil {
			em := fmt.Sprintf("internal error displaying Create page - no HTML template")
			log.Printf("%s\n", em)
			c.ErrorHandler(req, resp, em)
			return
		}
		err := page.Execute(resp.ResponseWriter, &form)
		if err != nil {
			em := fmt.Sprintf("Internal error while preparing create form after failed validation - %s",
				err.Error())
			log.Printf("%s\n", em)
			c.ErrorHandler(req, resp, em)
			return
		}
		return
	}

	// Create a person in the database using the validated data in the form
	dao := c.services.GetPeopleRepository()

	createdPerson, err := dao.Create(form.Person())
	if err != nil {
		// Failed to create person.  Display index page with error message.
		em := fmt.Sprintf("Could not create person %s - %s", form.Person().String(), err.Error())
		c.ErrorHandler(req, resp, em)
		return
	}

	// Success! Person created.  Display index page with confirmation notice
	notice := fmt.Sprintf("created new person %s", createdPerson.String())
	log.Printf("%s\n", notice)
	var listForm forms.ConcreteListForm
	listForm.SetNotice(notice)
	listPeople(req, resp, &listForm, c.services)
	return
}

// Edit fetches the data for the people record with the given ID and displays
// the edit page, populated with that data.
func (c Controller) Edit(req *restful.Request, resp *restful.Response,
	form forms.PersonForm) {

	log.SetPrefix("Edit() ")

	log.Println("parsing form")
	err := req.Request.ParseForm()
	if err != nil {
		// failed to parse form
		em := fmt.Sprintf("cannot parse form - %s", err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	// Get the ID of the person
	id := req.PathParameter("id")

	dao := c.services.GetPeopleRepository()
	// Get the existing data for the person
	person, err := dao.FindByIDStr(id)
	if err != nil {
		// No such person.  Display index page with error message.
		em := err.Error()
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	// Got the person with the given ID.  Put it into the form and validate it.
	// If the data is invalid, continue - the user may be trying to fix it.

	form.SetPerson(person)
	if !form.Validate() {
		em := fmt.Sprintf("invalid record in the people database - %s",
			person.String())
		log.Printf("%s\n", em)
	}

	// Display the edit page
	page := c.services.Template("Edit")
	if page == nil {
		em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	log.Printf("Execute\n")
	err = page.Execute(resp.ResponseWriter, form)
	if err != nil {
		// error while preparing edit page
		log.Printf("%s: error displaying edit page - %s", err.Error())
		em := fmt.Sprintf("error displaying page - %s", err.Error())
		c.ErrorHandler(req, resp, em)
	}
}

// Update responds to a PUT request.  For example:
// PUT /people/1
// It's invoked by the form displayed by a previous Edit request.  If the ID in the URI is
// valid and the request parameters from the form specify valid people data, it updates the
// record and displays the index page with a confirmation message, otherwise it displays
// the edit page again with the given data and some error messages.
func (c Controller) Update(req *restful.Request, resp *restful.Response,
	form forms.PersonForm) {

	log.SetPrefix("Update() ")

	// Get the person specified in the form from the DB.
	// (which also validates the id in the form).
	log.Printf("form=%v\n", form)
	if form.Person() == nil {
		em := fmt.Sprint("internal error - form should contain an updated person record")
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}

	dao := c.services.GetPeopleRepository()
	person, err := dao.FindByID(form.Person().ID())
	if err != nil {
		// There is no person with this ID.  The ID is chosen by the user from a
		// supplied list and it should always be valid, so there's something screwy
		// going on.  Display the index page with an error message.
		em := fmt.Sprintf("error searching for person with id %s - %s",
			form.Person().ID(), err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}

	// We have a matching person from the DB.
	log.Printf("got person %v\n", person)

	// Validate the new version of the person in the form.
	if !form.Validate() {
		// The data is invalid.  The validator has set error messages.  Return to
		// the edit screen.
		page := c.services.Template("Edit")
		if page == nil {
			em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
			log.Printf("%s\n", em)
			c.ErrorHandler(req, resp, em)
			return
		}
		err = page.Execute(resp.ResponseWriter, form)
		if err != nil {
			log.Printf("%s: error displaying edit page - %s", err.Error())
			em := fmt.Sprintf("error displaying page - %s", err.Error())
			c.ErrorHandler(req, resp, em)
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

		page := c.services.Template("Edit")
		if page == nil {
			em := fmt.Sprintf("internal error displaying Edit page - no HTML template")
			log.Printf("%s\n", em)
			c.ErrorHandler(req, resp, em)
			return
		}
		err = page.Execute(resp.ResponseWriter, form)
		if err != nil {
			// Error while recovering from another error.  This is looking like a habit!
			em := fmt.Sprintf("Internal error while preparing edit page after failing to update person in DB - %s", err.Error())
			log.Printf("%s\n", em)
			c.ErrorHandler(req, resp, em)
		} else {
			return
		}
	}

	// Success!  Display the index page with a confirmation notice
	notice := fmt.Sprintf("updated person %s", form.Person().String())
	log.Printf("%s:\n", notice)
	var listForm forms.ConcreteListForm
	listForm.SetNotice(notice)
	listPeople(req, resp, &listForm, c.services)
	return
}

// Delete reponds to a DELETE request and deletes the record with the given ID,
// eg DELETE http://server:port/people/1.
func (c Controller) Delete(req *restful.Request, resp *restful.Response) {

	log.SetPrefix("Delete()")

	err := req.Request.ParseForm()
	if err != nil {
		// failed - form does not parse
		em := fmt.Sprintf("Internal error - %s", err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	method := req.Request.FormValue("_method")
	if "DELETE" != method {
		// failed - _method param is not DELETE
		em := fmt.Sprintf("Internal error - request type %s must be DELETE", method)
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	id := req.PathParameter("id")

	dao := c.services.GetPeopleRepository()
	// Attempt the delete
	_, err = dao.DeleteByIDStr(id)
	if err != nil {
		// failed - cannot delete person
		em := fmt.Sprintf("Cannot delete person with id %s - %s", id, err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
		return
	}
	// Success - person deleted.  Display the index view with a notification.
	var form forms.ConcreteListForm
	notice := fmt.Sprintf("deleted person with ID %s", id)
	log.Printf("%s:\n", notice)
	form.SetNotice(notice)
	listPeople(req, resp, &form, c.services)
	return
}

// ErrorHandler displays the index page with an error message
func (c Controller) ErrorHandler(req *restful.Request, resp *restful.Response,
	errormessage string) {

	var form forms.ConcreteListForm
	form.SetErrorMessage(errormessage)
	listPeople(req, resp, &form, c.services)
}

// SetServices sets the services.
func (c *Controller) SetServices(services services.Services) {
	c.services = services
}

/*
 * The listPeople helper function fetches a list of people and displays the
 * index page.  It's used to fulfil an index request but the index page is
 * also used as the last page of a sequence of requests (for example new,
 * create, index).  If the sequence was successful, the form may contain a
 * confirmation note.  If the sequence failed, the form should contain an error
 * message.
 */
func listPeople(req *restful.Request, resp *restful.Response, form forms.ListForm,
	services services.Services) {

	log.SetPrefix("Controller.listPeople() ")

	dao := services.GetPeopleRepository()

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
	page := services.Template("Index")
	if page == nil {
		utilities.Dead(resp)
		return
	}
	err = page.Execute(resp.ResponseWriter, form)
	if err != nil {
		/*
		 * Error while displaying the index page.  We handle most internal
		 * errors by displaying the index page.  That's just failed, so
		 * fall back to the static error page.
		 */
		log.Printf(err.Error())
		page = services.Template("Error")
		if page == nil {
			utilities.Dead(resp)
			return
		}
		err = page.Execute(resp.ResponseWriter, form)
		if err != nil {
			// Can't display the static error page either.  Bale out.
			em := fmt.Sprintf("fatal error - failed to display error page for error %s\n", err.Error())
			log.Printf(em)
			panic(em)
		}
		return
	}
}
