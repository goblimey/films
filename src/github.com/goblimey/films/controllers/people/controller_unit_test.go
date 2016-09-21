package people

import (
	"errors"
	"fmt"
	// "html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
	peopleForms "github.com/goblimey/films/forms/people"
	mocks "github.com/goblimey/films/mocks/gomock"
	"github.com/goblimey/films/mocks/manual"
	pemocks "github.com/goblimey/films/mocks/pegomock"
	personModel "github.com/goblimey/films/models/person"
	retroTemplate "github.com/goblimey/films/retrofit/template"
	"github.com/goblimey/films/services"
	"github.com/golang/mock/gomock"
	"github.com/petergtz/pegomock"
)

var panicValue string

// TestUnitIndexWithOnePersonPE checks that PeopleHandler.Index() handles a list of
// people from FindAll() containing one person.  It uses pegomock to create a mock
// template.
func TestUnitIndexWithOnePersonPE(t *testing.T) {

	pegomock.RegisterMockTestingT(t)

	// Create a list containing one person.
	expectedID := uint64(42)
	expectedForename := "foo"
	expectedSurname := "bar"
	expectedPerson := personModel.MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	expectedPersonList := make([]personModel.Person, 1)
	expectedPersonList[0] = expectedPerson

	// Create the mocks and dummy objects.
	var url url.URL
	url.Opaque = "/people" // url.RequestURI() will return "/people"
	var httpRequest http.Request
	httpRequest.URL = &url
	httpRequest.Method = "GET"
	var request restful.Request
	request.Request = &httpRequest
	writer := pemocks.NewMockResponseWriter()
	var response restful.Response
	response.ResponseWriter = writer
	mockTemplate := pemocks.NewMockTemplate()
	var mockRepo = manual.MockRepo{expectedPersonList}
	page := make(map[string]retroTemplate.Template)
	page["Index"] = mockTemplate

	// Create a service that returns the mock repository and templates.
	var services services.ConcreteServices
	services.SetPeopleRepository(&mockRepo)
	services.SetTemplates(&page)

	// Create the form
	var form peopleForms.ConcreteListForm

	// The request supplies method "GET" and URI "/people".  Expect
	// template.Execute to be called and return nil (no error).
	pegomock.When(mockTemplate.Execute(writer, &form)).ThenReturn(nil)

	// Run the test.
	var controller Controller
	controller.SetServices(&services)
	controller.Index(&request, &response, &form)

	// We expect that the form contains the expected person list -
	// one person with id, forename and surname as expected.
	if form.People() == nil {
		t.Errorf("Expected a list, got nil")
	}

	if len(form.People()) != 1 {
		t.Errorf("Expected a list of 1, got %d", len(form.People()))
	}

	if form.People()[0].ID() != expectedID {
		t.Errorf("Expected ID %d, got %d",
			expectedID, form.People()[0].ID())
	}

	if form.People()[0].Forename() != expectedForename {
		t.Errorf("Expected forename %s, got %s",
			expectedForename, form.People()[0].Forename())
	}

	if form.People()[0].Surname() != expectedSurname {
		t.Errorf("Expected surname %s, got %s",
			expectedForename, form.People()[0].Surname())
	}

	// Check that the service were used as expected
	err := mockRepo.TestComplete()
	if err != nil {
		t.Errorf("TestUnitIndexWithOnePerson %s", err.Error())
	}
}

// TestUnitIndexWithOnePerson checks that PeopleHandler.Index() handles a list of
// people from FindAll() containing one person.  This is the same as the previous
// test, but uses gomock rather than pegomock.
func TestUnitIndexWithOnePerson(t *testing.T) {

	// Create a list containing one person.
	expectedID := uint64(42)
	expectedForename := "foo"
	expectedSurname := "bar"
	expectedPerson := personModel.MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	expectedPersonList := make([]personModel.Person, 1)
	expectedPersonList[0] = expectedPerson

	// Create the mocks and dummy objects.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var url url.URL
	url.Opaque = "/people" // url.RequestURI() will return "/people"
	var httpRequest http.Request
	httpRequest.URL = &url
	httpRequest.Method = "GET"
	var request restful.Request
	request.Request = &httpRequest
	mockWriter := mocks.NewMockResponseWriter(mockCtrl)
	var response restful.Response
	response.ResponseWriter = mockWriter
	mockTemplate := mocks.NewMockTemplate(mockCtrl)
	var mockRepo = manual.MockRepo{expectedPersonList}
	page := make(map[string]retroTemplate.Template)
	page["Index"] = mockTemplate

	// Create a service that returns the mock repository and templates.
	var services services.ConcreteServices
	services.SetPeopleRepository(&mockRepo)
	services.SetTemplates(&page)

	// Create the form
	var form peopleForms.ConcreteListForm

	// The request supplies method "GET" and URI "/people".  Expect
	// template.Execute to be called and return nil (no error).
	mockTemplate.EXPECT().Execute(mockWriter, &form).Return(nil)

	// Run the test.
	var controller Controller
	controller.SetServices(&services)
	controller.Index(&request, &response, &form)

	// We expect that the form contains the expected person list -
	// one person with id, forename and surname as expected.
	if form.People() == nil {
		t.Errorf("Expected a list, got nil")
	}

	if len(form.People()) != 1 {
		t.Errorf("Expected a list of 1, got %d", len(form.People()))
	}

	if form.People()[0].ID() != expectedID {
		t.Errorf("Expected ID %d, got %d",
			expectedID, form.People()[0].ID())
	}

	if form.People()[0].Forename() != expectedForename {
		t.Errorf("Expected forename %s, got %s",
			expectedForename, form.People()[0].Forename())
	}

	if form.People()[0].Surname() != expectedSurname {
		t.Errorf("Expected surname %s, got %s",
			expectedForename, form.People()[0].Surname())
	}

	// Check that the service were used as expected
	err := mockRepo.TestComplete()
	if err != nil {
		t.Errorf("TestUnitIndexWithOnePerson %s", err.Error())
	}
}

// TestUnitIndexWithErrorWhenFetchingPeoplePE checks that PeopleHandler.Index()
// handles errors from FindAll() correctly.  It uses pegomock to provide the
// mocked template.
func TestUnitIndexWithErrorWhenFetchingPeoplePE(t *testing.T) {

	log.SetPrefix("TestUnitIndexWithErrorWhenFetchingPeoplePE ")
	log.Printf("This test is expected to provoke error messages in the log")

	expectedErr := errors.New("Test Error Message")
	expectedErrorMessage := "error getting the list of people - Test Error Message"

	pegomock.RegisterMockTestingT(t)

	// Create the mocks and dummy objects.
	var url url.URL
	url.Opaque = "/people" // url.RequestURI() will return "/people"
	var httpRequest http.Request
	httpRequest.URL = &url
	httpRequest.Method = "GET"
	var request restful.Request
	request.Request = &httpRequest
	writer := pemocks.NewMockResponseWriter()
	var response restful.Response
	response.ResponseWriter = writer
	mockTemplate := pemocks.NewMockTemplate()
	var mockRepo manual.MockRepo
	page := make(map[string]retroTemplate.Template)
	page["Index"] = mockTemplate

	// Create a service that returns the mock repository and templates.
	var services services.ConcreteServices
	services.SetPeopleRepository(&mockRepo)
	services.SetTemplates(&page)

	// Create the form
	var form peopleForms.ConcreteListForm

	// The request supplies method "GET" and URI "/people".  Expect
	// template.Execute to be called and return the expected error.
	pegomock.When(mockTemplate.Execute(writer, &form)).ThenReturn(expectedErr)

	// Run the test.
	controller := MakeController(&services)

	controller.Index(&request, &response, &form)

	// Verify that the form contains the expected error message.
	if form.ErrorMessage() != expectedErrorMessage {
		t.Errorf("Expected error message to be %s actually %s", expectedErrorMessage, form.ErrorMessage())
	}

	err := mockRepo.TestComplete()
	if err != nil {
		t.Errorf("TestUnitIndexWithErrorWhenFetchingPeople %s", err.Error())
	}
}

// TestUnitIndexWithErrorWhenFetchingPeoplePE checks that PeopleHandler.Index()
// handles errors from FindAll() correctly.  This is the same as the previous test
// except that uses gomock to provide the mocked template.
func TestUnitIndexWithErrorWhenFetchingPeople(t *testing.T) {

	log.SetPrefix("TestUnitIndexWithErrorWhenFetchingPeople ")
	log.Printf("This test is expected to provoke error messages in the log")

	em := "Test Error Message"
	expectedErr := errors.New(em)
	expectedErrorMessage := fmt.Sprintf("error getting the list of people - %s", em)
	// Create a list containing one person.
	expectedID := uint64(42)
	expectedForename := "foo"
	expectedSurname := "bar"
	expectedPerson := personModel.MakeInitialisedPerson(expectedID, expectedForename, expectedSurname)
	expectedPersonList := make([]personModel.Person, 1)
	expectedPersonList[0] = expectedPerson

	// Create the mocks and dummy objects.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var url url.URL
	url.Opaque = "/people" // url.RequestURI() will return "/people"
	var httpRequest http.Request
	httpRequest.URL = &url
	httpRequest.Method = "GET"
	var request restful.Request
	request.Request = &httpRequest
	mockResponseWriter := mocks.NewMockResponseWriter(mockCtrl)
	var response restful.Response
	response.ResponseWriter = mockResponseWriter
	mockTemplate := mocks.NewMockTemplate(mockCtrl)
	var mockRepo = mocks.NewMockRepository(mockCtrl)
	page := make(map[string]retroTemplate.Template)
	page["Index"] = mockTemplate

	// Create a service that returns the mock repository and templates.
	var services services.ConcreteServices
	services.SetPeopleRepository(mockRepo)
	services.SetTemplates(&page)

	// Create the form
	var form peopleForms.ConcreteListForm

	// The request supplies method "GET" and URI "/people".  Expect
	// template.Execute to be called and return the expected error.  Expect
	// Execute to be called and return no error.  Expect the form to contain the
	// error message from FindAll and a nil list of people.
	mockRepo.EXPECT().FindAll().Return(nil, expectedErr)
	mockTemplate.EXPECT().Execute(mockResponseWriter, &form).Return(nil)

	// Run the test.
	controller := MakeController(&services)
	controller.Index(&request, &response, &form)

	// Verify that the form contains the expected error message.
	if form.ErrorMessage() != expectedErrorMessage {
		t.Errorf("Expected error message to be %s actually %s",
			expectedErrorMessage, form.ErrorMessage())
	}

	// Verify that the list of people is nil
	if form.People() != nil {
		t.Errorf("Expected the list of people to be nil.  Actually contains %d entries",
			len(form.People()))
	}
}

// TestUnitIndexWithManyFailures checks that PeopleHandler.Index() handles a series
// of errors correctly.
//
// Panic handling based on http://stackoverflow.com/questions/31595791/how-to-test-panics
//
func TestUnitIndexWithManyFailures(t *testing.T) {

	log.SetPrefix("TestUnitIndexWithManyFailures ")
	log.Printf("This test is expected to provoke error messages in the log")

	// Create the mocks and dummy objects.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	var url url.URL
	url.Opaque = "/people" // url.RequestURI() will return "/people"
	var httpRequest http.Request
	httpRequest.URL = &url
	httpRequest.Method = "GET"
	var request restful.Request
	request.Request = &httpRequest
	mockResponseWriter := mocks.NewMockResponseWriter(mockCtrl)
	var response restful.Response
	response.ResponseWriter = mockResponseWriter
	mockTemplate := mocks.NewMockTemplate(mockCtrl)
	mockErrorTemplate := mocks.NewMockTemplate(mockCtrl)
	var mockRepo = mocks.NewMockRepository(mockCtrl)

	// Create a template map containing the mock templates
	page := make(map[string]retroTemplate.Template)
	page["Index"] = mockTemplate
	page["Error"] = mockErrorTemplate

	// Create a service that returns the mock repository and templates.
	var services services.ConcreteServices
	services.SetPeopleRepository(mockRepo)
	services.SetTemplates(&page)

	var form peopleForms.ConcreteListForm

	// Expectations:
	// Index will run listPeople which will call repository.FindAll.  Make that
	// return an error, then listPeople will get the Index page from the template
	// and call its Execute method.  Make that fail, and the app will get the error
	// page and call its Execute method.  Make that fails and the app will panic
	// with a message "fatal error - failed to display error page for error ",
	// followed by the error message from the last Execute call.

	em1 := "first error message"
	expectedFirstErrorMessage := fmt.Sprintf(
		"some stuff - %s", em1)

	expectedFirstError := errors.New(em1)

	em2 := "second error message"
	expectedSecondError := errors.New(em2)

	em3 := "final error message"
	finalError := errors.New(em3)

	mockRepo.EXPECT().FindAll().Return(nil, expectedFirstError)
	// form will now be different (error message added) so don't compare it
	mockTemplate.EXPECT().Execute(mockResponseWriter, gomock.Any()).Return(expectedSecondError)
	mockErrorTemplate.EXPECT().Execute(mockResponseWriter, gomock.Any()).Return(finalError)

	// Expect a panic, catch it and check the value.  (If there is no panic,
	// this raises an error.)

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected the Index call to panic")
		} else {
			em := fmt.Sprintf("%v", r)
			// Verify that the panic value is as expected.
			if !strings.Contains(em, em3) {
				t.Errorf("Expected a panic with value containing \"%s\" actually \"%s\"",
					em3, em)
			}
		}
	}()

	// Run the test.
	controller := MakeController(&services)
	controller.Index(&request, &response, &form)

	// Verify that the form has an error message containing the expected text.
	if strings.Contains(form.ErrorMessage(), em1) {
		t.Errorf("Expected error message to be \"%s\" actually \"%s\"",
			expectedFirstErrorMessage, form.ErrorMessage())
	}

	// Verify that the list of people is nil
	if form.People() != nil {
		t.Errorf("Expected the list of people to be nil.  Actually contains %d entries",
			len(form.People()))
	}

}

// Recover from any panic and record the error.
func catchPanic() {
	log.SetPrefix("catchPanic ")
	if p := recover(); p != nil {
		em := fmt.Sprintf("%v", p)
		panicValue = em
		log.Printf(em)
	}
}
