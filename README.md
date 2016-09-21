# films
A simple MVC web server written in Go with restful interfaces providing CRUD operations on database tables.

The films web server is a very simple example of how the Go programming language can be used to build a website.  It has one table called "people" with three fields, a numeric ID, a text field "forename" and another "surname".  The server has web pages that perform the Create, Read, Update and Delete (CRUD) operations on the table.  It looks (and is) ridiculously simple, but this allows me to concentrate on the mechanics of handling web pages and data rather than the business logic that a real website would require.  If you are trying to figure out how to create a website or just do test-driven development in Go, looking at this project should give you a start.

The server is designed using the Model View Controller (MVC) pattern, which supports separation of concerns.  The major components of the web server are models, views and controllers.  Each of those components is only concerned with well-defined issues.  The models are repositories (sometimes called Data Access Objects or DAOs) that provide access to data in the database.  The models are not concerned with what will be done with the data or how it will be represented (rendered) when the user sees it.  The views are templates that produce HTML pages.  They are only concerned with rendering the data.  They make no decisions, they just render the data that they are given.  The controllers implement the business logic of the application.  They use the repositories to access the data and use the views to display it.  They are not concerned with where the data comes from or how it will be rendered.  They are only concerned with what the user is allowed to do with the data.  

I'm one of the organisers of the the Surrey Go User Group (https://groups.google.com/forum/#!forum/surrey-golang).  Having looked at this project, other members of the group are now writing their own web servers, borrowing ideas from this one as needed. 


Fetching and Building the Films Server
======================================

First, get the dependencies:

```
go get github.com/go-sql-driver/mysql
go get github.com/coopernurse/gorp
go get github.com/emicklei/go-restful
go get github.com/golang/mock/gomock
go get github.com/petergtz/pegomock/pegomock
```

Note: by default, go get does not update anything that you have already downloaded.  If you downloaded any of those packages a long time ago, you may wish to update them to the latest version using the -u flag, for example:

go get -u github.com/coopernurse/gorp

Next, clone the server source code and build it.  For example, in a UNIX or Linux command window:

```
cd $HOME
mkdir goprojects
cd goprojects
git clone https://github.com/goblimey/films
```

That will create a directory "goprojects" in your home directory containing a directory called "films".  

Now build the server:

```
cd films
. setenv.sh
go install github.com/goblimey/films
```

The only part of that which is specific to UNIX or Linux is running setenv.sh.  That adds the current directory to the GOPATH, and the bin directory to the PATH. (The "." at the start of that line is a command and is required.)

If you are running under Windows you can do the same by hand, something like this:

```
SET GOPATH=%GOPATH%;c:\Users\simon\goprojects\films
SET PATH=%PATH%;c:\Users\simon\goprojects\films\bin
```


Setting Up the Database
-----------------------

The server expects a MySQL database, so you need to install the MySQL client and server, which you can get from the Oracle website www.oracle.com.  You will need to create an account to download MySQL, but it is free.

Once you have MySQL running, create a database called "films". This must be accessible with all privileges by a user called "webuser" with password "secret".  You can set that up as follows using the mysql client:

```
$ mysql -u root -p
(supply the password for the MySQL root user)
mysql> create database films;
mysql> grant all on films.* to webuser identified by 'secret';
mysql> quit;
```

(If your MySQL server doesn't have a root password, omit the -p from the first command.)

Those login details (webuser/secret) are defined in the method MakeGorpMysqlDBSession() in utilities/dbsession/dbsession.go.  You can change them by editing that file and rebuilding the server.

Note that things like table and database names are case-sensitive when MySQL runs under UNIX, so the databases "FILMS", "Films" and "films" are different objects.  Under Windows those names would all apply to the same object.  (This is because the objects are represented by files and follow the naming rules for files on those systems.)

The server expects a table called "people".  if it doesn't exist, the server will create an empty one when you start it up.  If you prefer to set one up yourself, here is a suitable description:

```
mysql> use films;
mysql> describe people;
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | mediumint(9) | NO   | PRI | NULL    | auto_increment |
| forename | varchar(100) | NO  |     | NULL    |                |
| surname  | varchar(100) | NO  |     | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
3 rows in set (0.00 sec)
```


Running the Server
------------------

Start the server from a command window.  If you have logged out and back in again since building the server, or if you are doing this part in a new command window, move to the same directory as before ($HOME/goprojects/films in my example) and set up your GOPATH and PATH variables.  Under UNIX and Linux, that's:

```
    .  setenv.sh
```

The "go install" that you ran earlier created a program called "films" in the bin directory at the top level of the project.  The commands in setenv.sh put the program into your path, so you can run it.

For Windows see above.

Move to the directory containing the views directory:

```
    cd src/github.com/goblimey/films
```

Run the server:

```
     films
```

The server listens on port 4000.  In a web browser, navigate to

    http://localhost:4000/people

Initially there will be no entries in the people table.  Use the create button to create some.

The create screen has some simple validation to ensure that you fill in both fields.  Try missing one or both of them out and pressing the submit button.

To stop the web server, go to the command window from which it is being run, hold down the ctrl key and type a single "c".  The result is instant, you don't need to hit the enter key.


How the Server Works
====================

The server is composed of a main program, models, views and controllers.  The controllers stand between the models and the views and respond to HTML request made by the user.  For each request, a controller runs, gets data from a model, manipulates it according to the business logic and then sends the result to a view which renders it.

I've borrowed a few ideas from the Java MVC frameworks Struts and Spring.

Struts has a useful gadget called a form bean, which is a data transfer object used to carry information around the controllers and views.  I have a object called "PersonForm" that holds information about a single person and another called ListForm that holds a list of people for pages such as the Index page that display lists of people.

Spring implements Inversion of Control to support interfaces and testing using mocking.  Each method in a Java Spring controller takes a standard set of arguments and is separated from the details of request handling.  The arguments are specified in terms of interfaces.  Something outside the controller creates the objects, so the controller doesn't know what type they are, it just knows that they satisfy the interfaces.

The films server follows a simple version of those ideas.  Each method in a controller implements a request.  All incoming requests are fielded by the main.marshall method.  This figures out which controller to use and which method to call to handle the request.  The controller methods take a standard set of arguments defined using interfaces rather than real structures.  The database repositories that supply data to the model also supply it as interfaces.  This allows objects that conform to the same interface to be used interchangeably.  In particular mocks and dummies can be used during testing.

I've created objects modelled on the Struts form beans to carry the data to be rendered by the views.  Each form object can contain a notification message and/or an error message to be displayed at the top of the page, data items to be displayed and error messages about the data items.  When the user submits data in an HTML form, the main.marshall method presents it to the controller method in a form object.  When the server executes a template to display a web page, it supplies the data for the page in a form object.

The standard Go library includes a library net/html, which provides a framework for building and displaying web pages.  I use this to provide the views.  Each view takes a package of data provided by the controller, creates an HTML page to display it and sends the page to the user's browser.  The contents of the view is determined by a form object.

For example, this interface defines the form object used to carry data about a Person:

```go
type PersonForm interface {
	// Person gets the Person embedded in the form.
	Person() personModel.Person
	// Notice gets the notice.
	Notice() string
	// ErrorMessage gets the general error message.
	ErrorMessage() string
	// FieldErrors returns all the field errors as a map.
	FieldErrors() map[string]string
	// ErrorForField returns the error message about a field (may be an empty string).
	ErrorForField(key string) string
	// String returns a string version of the PersonForm.
	String() string
	// SetPerson sets the Person in the form.
	SetPerson(person personModel.Person)
	// SetNotice sets the notice.
	SetNotice(notice string)
	//SetErrorMessage sets the general error message.
	SetErrorMessage(errorMessage string)
	// SetErrorMessageForField sets the error message for a named field
	SetErrorMessageForField(fieldname, errormessage string)
	// Validate validates the data in the Person and sets the various error messages.
	// It returns true if the data is valid, false if there are errors.
	Validate() bool
}
```

A structure that satisfies this interface contains a Person record with fields ID, surname and forename, a notice that should be displayed at the top of the page, an error message that should be displayed at the top of the page, and a list of field errors - error messages about individual fields of the Person record.  It provides Validate method which the server uses to check incoming data.

In the films server, all requests are sent to a function in the main package called marshal.  This figures out which controller to call (currently there's only one, the people controller).  It creates a controller and a form object containing the data that the controller needs.  It calls a method in the controller to handle the request, supplying the data as an argument.  The handler does some work and then uses a template to render an HTML page displaying the result.

A quick explanation for anybody who hasn't written a back-end web server before:  A web server loops forever, waiting for a request and then calling a controller method to handle it.  The controller method sends a response back to the browser which is an HTML page.  For the user's session to continue, every response page should contain buttons or links that allow them to issue another request.  The user's session is a series of requests and responses which lasts until the user gets bored and goes away.  The server runs until it's forcibly shut down.

So the user and the server run through a session composed of a series of requests and responses.  Every response displays an HTML page and every page has links or buttons to issue another request and continue the session.  The pages allow the user to create, read, update and delete data in the table.  They have no direct access to the database or its tables, so they can only do with this table what the controller allows them to do and they can't access any other tables, if any exist.

The workings of the server are best illustrated with an example.  A user starts at the index page for the people controller http://localhost:4000/people.  This has a button that issues a request to create a Person, and the user presses it.  The browser sends the request to the server, which runs main.marshall to field it.  This creates a people controller and an empty person form and calls the controller's New method, passing the form.  New executes the Create template, passing the empty form as data.  The user sees a web form with empty fields and a submit button.

The user fills in the surname but misses out the forename, then presses the submit button.  Both surname and forename are mandatory fields, so this request should be rejected.  The form sends a Create request which is fielded by main.marshal.  It reads the form data, sets up a new person form containing the supplied surname, creates a new people controller and calls its Create method.  This runs the form's Validate method, which rejects the forename field.  The server adds a message to the form in the field errors list for the forename field, then executes the Create template, passing it the form.  The user sees a create page again, but this time the surname field contains the text that they supplied and there is an error message next to the forename field.

Whenever a user submits a request, the server may hit a problem which is not to do with a particular field, for example, it cannot connect to the database.  In that case, it creates a form with the ErrorMessage field set and passes that to a template.  The user sees a page with the error message at the top.

The user fills in the forename and submits the form again.  This time validation is successful.  The people server's create method creates a new person in the database.  Now we want the user to see a page with a notification at the top saying that a new Person record was created.  We show them the people resource's index page which also lists all the people records.  The Index page is driven by a ListForm, not a PersonForm.  The server gets a list of person records from the database, creates a ListForm containing that list, adds the notification and executes the Index template to display the index page page including the notification.  Notifications are displayed in green.  (Errors are displayed in red.)

The index page contains buttons and links that allow the user to send a request to either: view the details of one of the person records in the list; edit a record; delete a record; create a new record.

The requests and web pages are laid out using the REST model, implemented using the go-restful library.  A RESTful web server provides a set of resources that the user can access.  Each resource has its own model and controller, plus a set of views.  A resource can be (but need not be) represented by a database table.  All requests concerning a resource follow a pattern that starts with the resource name, for example:

```
    /people             display all people
    /people/42          display the person with ID 42
    /people/42/edit     fetch the data for person 42 and display a screen to change it
```

This server is called films because a future version will display information about films - a very simple form of IMDb.  The people table will hold data about actors, directors and so on, and there will be other tables, with web pages to manipulate them.  At present there is one resource, representing one table, so the server has one model, one controller and one set of views.

In Go, an object satisfies an interfaces if it implements all of the methods of the interface, so you can create an interface and then create an object that satisfies it, or you can take an existing object and write an interface that it satisfies.  Unlike with languages such as Java, this means that an interface can be retrofitted to an object that you did not create.  For example, the net/html package defines a structure called a Template which is used to render an HTML page.  The films server includes an interface Template that the HTML template structure satisfies.  The controller uses the templates in terms of that interface, so its templates can be replaced with mock versions during testing.

I've made extensive use of factory functions to create objects.  For example, I have an interface Person and a structure ConcretePerson that satifies it.  ConcretePerson has this function:

```go
// MakePerson creates and returns a new uninitialised Person object
func MakePerson() Person {
    var concretePerson ConcretePerson
    return &concretePerson
}
```

This is a function, not a class method, so it's called like so:
```go
var person Person = ConcretePerson.MakePerson()
``` 

ConcretePerson is a structure, not an instance of a structure, so MakePerson is a function, not a method.  It can only be called via the name of the structure as above. It's important that it's not a method, because you can only call a method on an instance of a structure and we don't have one of those yet - we're calling the factory function to create one.  (In Java, the equivalent of MakePerson is called a static factory method.)

MakePerson creates an empty ConcretePerson but returns it as a Person.

A Go interface can only define class methods, so you can't wrote an interface that represents the MakePerson function.

Having called the factory function, we have an object called person which is defined by the Person interface (not by the ConcretePerson structure).  Its class methods include some setters, so we can put data into the object:

```go
person.SetSurname("Simon")
person.SetForeName("Ritchie")
```

If you pass this object to a method, the method only knows that it satisfies the Person interface.  So one piece of software (usually main.marshall) can create an object using a factory function and pass it to a controller method to do the work.  The general rule is that the stuff that does the work doesn't know or care what the object is that it's working on, or how it was created.  It just knows which interface it satisfies.  This makes it easy to test the controller methods thoroughly.

(At present, some controller methods also call the factory functions, but they only do that so that they can call another controller method and pass the object to it as an interface - if something uses a factory function to create an object, it promptly passes it to something else to do work on it.  The people controller's errorHandler function does this.  This is unfortunate, because it means that the controller has to be polluted with knowledge of the real objects that are being used.  It would be better if it could be written to only use stuff that's passed into it.

All of the basic objects in this project (Person, PersonForm, ListForm and so on) are defined by interfaces and for each interface there is a concrete structure that satisfies the interface and provides factory functions to create an object of that type, returning it as an interface. 

I've also created a services object, which provides functionality that all controllers need.  When main.marshall creates an instance of a controller, it binds the services object into it.  The services object supplies the HTML templates and the repository classes that give access to the database tables.  Again these are defined in terms of interfaces, so during testing, a dummy version of the services can be substituted. 

(An obvious solution to my pollution issue is to use the services layer to provide the factory methods, but that's harder than it looks.  My first attempt led to circular dependencies, where class A includes class B and class includes class A.  That's not allowed in Go.)


The Database
============

The database runs under MySQL.  The server uses the GORP library to connect to the database.  That, of course, is the concern only of the model, and another model can be slotted in to replace it.  In the unit tests, the MySQL model is replaced by objects that provide fixed datasets that drive the logic through the desired logic paths.

Testing
=======

The films server includes unit tests, each of which tests a single unit of software in isolation by providing it with dummy objects containing data specifically written for the test.  There are also integration tests, where a few units of software are bound together and tests are run to check that everything hangs together.  My integration tests check that the controller and the MySQL model work together correctly.  Finally, system tests check that the whole system works together.  Go includes a facility for running tests on a complete web server, but I use Selenium, specifically the Firefox Selenium addon, which can be used to test any web server, regardless of the technology used to create it.

The test directory contains a number of Selenium scripts for system testing.  To run them in Firefox, install the Selenium addon, start it up and use the file menu to load one of the test suites.  They are in the tests directory.  (Don't forget to start the films server!)  Use the green arrow buttons to run the whole suite or one of the tests in it.  As the test runs, you can see the results in your browser.  There's also a control to speed up and slow down the replay, so you can watch what's going on.  The tests assume that the database is empty at the start.

The unit and integration tests are written using the standard Go test facilities, plus gomock and pegomock for mocking.  These tests live in the same directory as the module that they are testing, so it's easy to see how it's been tested.

Go offers a naming convention that can be used to classify tests.  The test controller can be made to run only those tests whose names match a pattern.  I use this to differentiate between unit and integration tests.  Each method of a unit test has a name that starts "TestUnit", for example "TestUnitCreatePersonFormAndRetrievePerson".  Each integration test has a name that starts "TestInt".

On a UNIX or Linux server, the tests can be run using the shell script test.sh as follows:

```
    ./test.s unit
```

runs just the unit tests

```
    ./test.sh int
```

runs just the integration tests

```
    ./test.sh
```

runs both the unit and integration tests.



Mocking
=======

At present I'm experimenting with using hand-crafted mocks and also with generating them automatically using the mocking frameworks gomock and pegomock.  I have tests that use all three techniques.

Neither gomock nor pegomock are awfully well-documented.  In particular, the documentation assumes that you have already used another mocking framework and understand how they work.  In case you don't, the fundamental idea is that, given an interface, the mocking framework produces a structure (a mock) that satisfies the interface, and which you can control.

Your test should call some method, passing it objects that drive it through a particular path of logic.  First your test should set expectations, which means that it should configure the mock to return the right values in the right sequence to drive the method under test.    

For example, the people controller package (github.com/goblimey/films/controllers/people) contains a test program that runs a series of unit tests using mock version of the HTML template, including TestUnitIndexWithOnePersonPE.  At the start of that test, I create a mock template and a services object that will supply it the controller:

```go
    mockTemplate := pemocks.NewMockTemplate()
    page := make(map[string]retroTemplate.Template)
    page["Index"] = mockTemplate
    var services services.ConcreteServices
    services.SetTemplates(&page)
```

Next, I set the mock's expectations:

```go
    pegomock.When(mockTemplate.Execute(writer, &form)).ThenReturn(nil)
```

That means "When the mock's Execute method is called with those arguments, return the value nil".  In this case, the arguments are objects that the test has set up and configured.
 
The method under test in this case is the people controller's index method.  Next I create a controller, pass in my dummy services object, and call the method:

```go
    // Run the test.
    var controller Controller
    controller.SetServices(&services)
    controller.Index(&request, &response, &form)
```

Looking at the controller's source code, Index just calls listPeople, which has this code at the end:

```go
    // Display the index page
    page := services.Template("Index")
    if page == nil {
        utilities.Dead(resp)
        return
    }
    err = page.Execute(resp.ResponseWriter, form)
    if err != nil {
```

This gets the template for the Index page from the services object, checks that it's not a nil pointer, then calls its Execute method.  Execute returns an error object, nil if there is no error.  If Execute doesn't return an error, listPeople returns, job done.

In our test, the services object is a dummy that supplies the mock template, so the controller calls the mock's Execute method.  We configured the mock to return nil, so the controller heads down the logic path that it would follow in the real world if it called a real HTML template's Execute method and got nil as the return value.

(If you wanted to test what the controller would do if it got back an error, then when you set up the expectations, you would create an error object and tell the mock to return that instead of nil.  The test program has other tests that do that.)
 
At the end of the test, the mock is torn down and as part of that sequence it automatically checks that everything went according to plan.  In particular it checks that its method were called in the sequence that you defined in the expectations, and raises errors if they weren't.

You can also add your own checks to the test, for example after you have called the controller method, the test can examine the contents of the form that you passed.

I don't include the mocks in my git repository.  This is deliberate, as they are generated by my test script test.sh, and that is in the repository.  (If you are running under Windows, the script won't work and you will have to generate the mocks yourself, but you can use it as a guide.)

There are other mocking frameworks available.  The pegomock github page includes a survey of them. 

It's worth saying that not everybody in the Go community agrees with the idea of using a mocking framework for test-driven development.  I find them a useful and cost-effective way to test my software.

Go also offers a framework that you can use to build a complete set of system tests for a web server, making web requests and checking that the correct response comes back.  You can use that, but of course, you still have to test your server using a web browser, preferably using all the common web browsers, to make sure that the pages look sensible.  I use Selenium for system testing web servers because it's visual, it can record tests interactively and it runs in a real browser, so I can see the pages as the tests run.
