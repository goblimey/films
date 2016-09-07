package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	htmlTemplate "html/template"

	restful "github.com/emicklei/go-restful"
	peopleController "github.com/goblimey/films/controllers/people"
	retrofitrequest "github.com/goblimey/films/retrofit/restful/request"
	retrofitresponse "github.com/goblimey/films/retrofit/restful/response"
	"github.com/goblimey/films/template"
)

// templates is a map of maps to serve out the HTML templates.  It has one
// subsidiary map for each controller.
var templateMap map[string]map[string]template.Template

// peopleRequestRE is the regular expression for the URI of a request to be handled
// by the PeopleController - for example: "/people" or "/people/1/delete".
var peopleRequestRE = regexp.MustCompile(`^/people$|^/people/.*`)

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

		} else {
			// some other error
			log.Println(err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
		}

		os.Exit(-1)
	}

	if !fileInfo.IsDir() {
		// views exists but is not a directory
		em := "the file views must be a directory"
		log.Println(em)
		fmt.Fprintln(os.Stderr, em)
	}

	// Set up the map of templates.  If anything goes wrong, report the error,
	// but at least we now think that we are in the right directory.
	templateMap = setupTemplates()

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

// setupTemplates creates a map of maps to serve out the templates, one major
// entry per controller.  If anyting goes wrong, the Must call will panic.
func setupTemplates() map[string]map[string]template.Template {

	templates := make(map[string]map[string]template.Template)

	// Set up the templates for the people controller.
	templates["people"] = make(map[string]template.Template)

	// This is the template for the error page, shared by all controllers.
	tp := htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/html/error.html"))

	var errorPageTP template.ConcreteTemplate
	errorPageTP.SetHTMLTemplate(tp)
	templates["people"]["Error"] = &errorPageTP

	tp = htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/index.ghtml",
	))

	var peopleIndexTP template.ConcreteTemplate
	peopleIndexTP.SetHTMLTemplate(tp)

	templates["people"]["Index"] = &peopleIndexTP

	tp = htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/create.ghtml",
	))

	var peopleCreateTP template.ConcreteTemplate
	peopleCreateTP.SetHTMLTemplate(tp)
	templates["people"]["Create"] = &peopleCreateTP

	tp = htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/show.ghtml",
	))

	var peopleShowTP template.ConcreteTemplate
	peopleShowTP.SetHTMLTemplate(tp)
	templates["people"]["Show"] = &peopleShowTP

	tp = htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/edit.ghtml",
	))

	tp = htmlTemplate.Must(htmlTemplate.ParseFiles(
		"views/templates/_base.ghtml",
		"views/templates/people/edit.ghtml",
	))
	var peopleEditTP template.ConcreteTemplate
	peopleEditTP.SetHTMLTemplate(tp)
	templates["people"]["Edit"] = &peopleEditTP

	return templates
}

// marshall passes the request to the appropriate controller.
func marshall(request *restful.Request, response *restful.Response) {

	log.SetPrefix("main.marshall() ")

	defer catchPanic()

	uri := request.Request.URL.RequestURI()

	log.Println("uri=", uri)

	// Create a mockable request and response
	mockableRequest := retrofitrequest.MakeRequest(request)
	mockableResponse := retrofitresponse.MakeResponse(response)

	if peopleRequestRE.MatchString(uri) {
		log.Printf("Sending request %s to PeopleController\n",
			mockableRequest.HTTPRequest().URL().RequestURI())
		if templateMap["people"] == nil {
			log.Println("baling out - no templates for people controller")
			return
		}
		var controller peopleController.Controller
		controller.Marshall(mockableRequest, mockableResponse, templateMap["people"])
	}
}

// Recover from any panic and log an error.
func catchPanic() {
	if p := recover(); p != nil {
		log.Printf("unrecoverable internal error %v\n", p)
	}
}
