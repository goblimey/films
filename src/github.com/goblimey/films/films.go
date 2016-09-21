package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	restful "github.com/emicklei/go-restful"
	peopleController "github.com/goblimey/films/controllers/people"
	forms "github.com/goblimey/films/forms/people"
	personModel "github.com/goblimey/films/models/person/gorpmysql"
	peopleRepo "github.com/goblimey/films/repositories/people"
	retroTemplate "github.com/goblimey/films/retrofit/template"
	"github.com/goblimey/films/services"
	"github.com/goblimey/films/utilities/dbsession"
)

// peopleRequestRE is the regular expression for the URI of any request to be
// handled by the PeopleController - for example: "/people", "/people/1/delete"
// and so on.
var peopleRequestRE = regexp.MustCompile(`^/people$|^/people/.*`)

// The following regular expressions are for specific request URIs, to work out
// which controller method to call.  For example, a GET request with URI "/people"
// produces a call to the Index method of the people controller.
//
// The requests follow the REST model and therefore carry data such as IDs
// in the request URI rather than in HTTP parameters, for example:
//
//    GET /people/435
//
// rather than
//
//    GET/people&id=435
//
// Only form data is supplied through HTTP parameters

// The peopleDeleteRequestRE is the regular expression for the URI of a delete
// request containing a numeric ID - for example: "/people/1/delete".
var peopleDeleteRequestRE = regexp.MustCompile(`^/people/[0-9]+/delete$`)

// The peopleShowRequestRE is the regular expression for the URI of a show
// request containing a numeric ID - for example: "/people/1".
var peopleShowRequestRE = regexp.MustCompile(`^/people/[0-9]+$`)

// The peopleEditRequestRE is the regular expression for the URI of an edit
// request containing a numeric ID - for example: "/people/1/edit".
var peopleEditRequestRE = regexp.MustCompile(`^/people/[0-9]+/edit$`)

// The peopleUpdateRequestRE is the regular expression for the URI of an update
// request containing a numeric ID - for example: "/people/1".  The URI
// is the same as for the show request, but we give it a different name for
// clarity.
var peopleUpdateRequestRE = peopleShowRequestRE

// page is a map of html templates, the views for the people resource.
var page *map[string]retroTemplate.Template

func main() {
	log.SetPrefix("main() ")
	log.Println("startup")

	// Nothing is going to work without the templates in the views directory.
	// If there is no views directory, give up.  Most likely, the user has not
	// moved to the right directory before running this.
	fileInfo, err := os.Stat("views")
	if err != nil {
		if os.IsNotExist(err) {
			// views does not exist
			em := "cannot find the views directory"
			log.Println(em)
			fmt.Fprintln(os.Stderr, em)

		} else if !fileInfo.IsDir() {
			// views exists but is not a directory
			em := "the file views must be a directory"
			log.Println(em)
			fmt.Fprintln(os.Stderr, em)

		} else {
			// some other error
			log.Println(err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
		}

		os.Exit(-1)
	}

	// Set up the map of templates.
	page = createPeopleTemplates()

	// Set up the restful web service.  Send all requests to marshall().

	ws := new(restful.WebService)
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("views/stylesheets"))))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("views/html"))))

	// Tie all expected requests to the marshall.
	ws.Route(ws.GET("/people").To(marshall))
	ws.Route(ws.GET("/people/{id}/edit").To(marshall))
	ws.Route(ws.GET("/people/{id}").To(marshall))
	ws.Route(ws.GET("/people/create").To(marshall))
	ws.Route(ws.POST("/people").Consumes("application/x-www-form-urlencoded").To(marshall))
	ws.Route(ws.POST("/people/{id}").Consumes("application/x-www-form-urlencoded").To(marshall))
	ws.Route(ws.POST("/people/{id}/delete").Consumes("application/x-www-form-urlencoded").To(marshall))
	restful.Add(ws)

	log.Println("starting the listener")
	err = http.ListenAndServe(":4000", nil)
	log.Println("baling out - " + err.Error())
}

// createPeopleTemplates creates a map to serve out the templates for the people
// controller.  If anything goes wrong, the Must call will panic.  The function
// assumes that the current directory contains a views directory containing the
// views.
func createPeopleTemplates() *map[string]retroTemplate.Template {

	templates := make(map[string]retroTemplate.Template)

	// This is the template for the error page, shared by all controllers.
	errorTP := template.Must(template.ParseFiles(
		"views/html/error.html"))

	templates["Error"] = errorTP

	peopleIndexTP := template.Must(template.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/index.ghtml",
	))

	// var peopleIndexTP template.ConcreteTemplate
	// peopleIndexTP.SetHTMLTemplate(tp)

	templates["Index"] = peopleIndexTP

	peopleCreateTP := template.Must(template.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/create.ghtml",
	))

	//var peopleCreateTP template.ConcreteTemplate
	//peopleCreateTP.SetHTMLTemplate(tp)
	templates["Create"] = peopleCreateTP

	peopleShowTP := template.Must(template.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/show.ghtml",
	))

	//var peopleShowTP template.ConcreteTemplate
	//peopleShowTP.SetHTMLTemplate(tp)
	templates["Show"] = peopleShowTP

	peopleEditTP := template.Must(template.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/edit.ghtml",
	))

	//var peopleEditTP template.ConcreteTemplate
	//peopleEditTP.SetHTMLTemplate(tp)
	templates["Edit"] = peopleEditTP

	return &templates
}

// marshall passes the request and response to the appropriate method of the
// appropriate  controller.
func marshall(request *restful.Request, response *restful.Response) {

	log.SetPrefix("main.marshall() ")

	defer catchPanic()

	// Create a service supplier
	session, err := dbsession.MakeGorpMysqlDBSession()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	var repo peopleRepo.GorpMysqlRepo
	repo.SetSession(session)
	var services services.ConcreteServices
	services.SetPeopleRepository(&repo)
	services.SetTemplates(page)

	uri := request.Request.URL.RequestURI()

	log.Println("uri=", uri)

	// The REST model uses HTTP requests such as PUT and DELETE.  The standard browsers do not support
	// these operations, so they are implemented using a POST request with a parameter "_method"
	// defining the operation.  (A post with a parameter "_method=PUT" simulates a PUT, and so on.)

	method := request.Request.Method
	if method == "POST" {
		// handle simulated PUT, DELETE etc via the _method parameter
		simMethod := request.Request.FormValue("_method")
		if simMethod == "PUT" || simMethod == "DELETE" {
			method = simMethod
		}
	}
	log.Printf("method %s", method)

	if peopleRequestRE.MatchString(uri) {

		log.Printf("Sending request %s to PeopleController\n", uri)

		var controller = peopleController.MakeController(&services)

		// Call the appropriate handler for the request

		switch method {

		case "GET":

			if uri == "/people" {
				// "GET http://server:port/people" - fetch all the valid people
				// records and display them.
				var form forms.ConcreteListForm
				controller.Index(request, response, &form)

			} else if peopleEditRequestRE.MatchString(uri) {

				// "GET http://server:port/people/1/edit" - fetch the people record
				// given by the ID in the request and display the form to edit it.
				var form forms.ConcretePersonForm
				controller.Edit(request, response, &form)

			} else if uri == "/people/create" {

				// "GET http://server:port/people/create" - display the form to
				// create a new people record.
				var form forms.ConcretePersonForm
				// Create an empty person to get started.
				person := personModel.MakePerson()
				form.SetPerson(person)
				controller.New(request, response, &form)

			} else if peopleShowRequestRE.MatchString(uri) {

				// "GET http://server:port/people/435" - fetch the people record
				// with ID 435 and display it.

				// Pass the ID to the controller via the form - get the ID from
				// the request, create a person containing (just) that ID, put that
				// person into the form.
				var form forms.ConcretePersonForm
				idStr := request.PathParameter("id")
				log.Printf("show id=%s", idStr)
				id, err := strconv.ParseUint(idStr, 10, 64)
				if err != nil {
					// This request is normally made from a link in a view.  The
					// link should always be correct, so this should never happen!
					em := fmt.Sprintf("illegal id %s", idStr)
					log.Println(em)
					controller.ErrorHandler(request, response, em)
				}
				person := personModel.MakePerson()
				person.SetID(id)
				form.SetPerson(person)
				controller.Show(request, response, &form)
			}

		case "PUT":
			if peopleUpdateRequestRE.MatchString(uri) {

				// POST http://server:port/people/1" - update the people record with
				// the given ID from the URI using the form data in the body.
				form := getPersonFormFromRequest(request, response, controller,
					&services)
				controller.Update(request, response, form)

			} else if uri == "/people" {

				// POST http://server:port/people" - create a new people record from
				// the form data in the body.
				form := getPersonFormFromRequest(request, response, controller,
					&services)
				controller.Create(request, response, form)
			}

		case "DELETE":
			if peopleDeleteRequestRE.MatchString(uri) {

				// "POST http://server:port/people/1/delete" - delete the people
				// record with the ID given in the request.
				controller.Delete(request, response)
			}

		default:
			em := fmt.Sprintf("unexpected HTTP method %v", method)
			log.Println(em)
			controller.ErrorHandler(request, response, em)
		}
	}
}

// getPersonFormFromRequest gets the person data from the request, creates a
// GorpMySQLPerson and returns it in a PersonForm.
func getPersonFormFromRequest(req *restful.Request, resp *restful.Response,
	c peopleController.Controller, services services.Services) forms.PersonForm {

	log.SetPrefix("getPersonFormFromRequest() ")

	err := req.Request.ParseForm()
	if err != nil {
		em := fmt.Sprintf("cannot parse form - %s", err.Error())
		log.Printf("%s\n", em)
		c.ErrorHandler(req, resp, em)
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
			c.ErrorHandler(req, resp, em)
			return nil
		}
		person.SetID(id)
	}
	person.SetForename(strings.TrimSpace(req.Request.FormValue("forename")))
	person.SetSurname(strings.TrimSpace(req.Request.FormValue("surname")))
	form.SetPerson(&person)
	log.Printf("form %s\n", form.String())
	return &form
}

// Recover from any panic and log an error.
func catchPanic() {
	if p := recover(); p != nil {
		log.Printf("unrecoverable internal error %v\n", p)
	}
}
