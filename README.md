# films
A simple MVC web server written in Go with restful interfaces providing CRUD operations on database tables.

The films web server is a very simple example of how the Go programming language can be used to build a website.  It has one table called "people" with three fields, a numeric ID, a text field "forename" and another "surname".  The server has web pages that perform the Create, Read, Update and Delete (CRUD) operations on the table.  It looks (and is) ridiculously simple, but this allows me to concentrate on the mechanics of handling web pages and data rather than the business logic that a real websiote would require.  If you are trying to figure out how to create a website or just do test-driven development in Go, this software is intended to give you a good start.

I'm one of the organisers of the the Surrey Go User Group.  (https://groups.google.com/forum/#!forum/surrey-golang)  This software is the result of a piece of work that I'm doing with that group.  

The server is designed using the Model View Controller (MVC) pattern, which supports separation of concerns.  The major components of the web server are models, views and controllers.  Each of those components is only concerned about well-defined issues.  The models represent the data in the database, along with a bit of software called a repository which gives access to it.  The model is not concerned with what will be done with the data or how it will be represented (rendered) when the user sees it.  The views are HTML pages.  They are only concerned with rendering the data.  They make no decisions, they just render what they are given.  The controllers access data in the databases and use the views to display it.  The controllers are not concerned with where the data comes from or how it will be rendered.  They are only concerned with what the user is allowed to do with the data.  They implement the business logic of the application.

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

The only part of that which is specific to UNIX or Linux is running setenv.sh.  That adds the current directory to the GOPATH, and the bin directory to the PATH.  

If you are running under Windows you can do the same by hand, something like this:

```
SET GOPATH=%GOPATH%;c:\Users\simon\goprojects\films
SET PATH=%PATH%;c:\Users\simon\goprojects\films\bin
```

Setting Up the Database
-----------------------

The server expects a MySQL database, so you need to install the MySQL client and server, which you can get from the Oracle website www.oracle.com.

Once you have MySQL running, create a database called "films". This must be accessible with all privileges by a user called "webuser" with password "secret".  You can set that up as follows using the mysql client:

```
$ mysql -u root -p
(supply the password for the MySQL root user)
mysql> create database films;
mysql> grant all on films.* to webuser identified by 'secret';
mysql> quit;
```

(If your MySQL server doesn't have a root password, omit the -p from the first command.)

Those login details (webuser/secret) are defined in the method MakeGorpMysqlDBSession() in utilities/dbsession/dbsession.go.  you can change them by editing that file and rebuilding the server.

Note that things like table and database names are case-sensitive when MySQL runs under UNIX, so the databases "FILMS", "Films" and "films" are different objects.  Under Windows those names would all apply to the same object.  (This is because the objects are represented by files and follow the naming rules for files on those systems.)

The server expects a table called "people".  if it doesn't exist, you can create an empty one.  If you prefer to set one up, here is a suitable description:

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

Start the server from a command window as follows:

If you haven't already, in the top-level directory, run setenv.sh like so:

```
    .  setenv.sh
```

(The "." at the start is a command and must be present.)

That puts the program "films" into your path.

Change to the directory containing the views directory:

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


How the Server Works
====================

The server is composed of a main program, models, views and controllers.  The controllers stand between the models and the views and respond to HTML request made by the user.  For each request, a controller runs, gets data from a model, manipulates it according to the business logic and then sends the result to a view which renders it.

I've borrowed a few ideas from the Java MVC frameworks Struts and Spring.

Struts has a useful gadget called a form bean, which is a data transfer object used to carry information around the controllers and views.  I have a object called "PersonForm" that holds information about a single person and another called ListForm that holds a list of people for pages such as the Index page that display lists of people.

Spring implements Inversion of Control to support interfaces and testing using mocking.  Each method in a Java Spring controller takes a standard set of arguments and is separated from the details of request handling.  The arguments are specified in terms of interfaces.  Something outside the controller creates the objects, so the controller doesn't know what type they are, it just knows that they satisfy the interfaces.  

The films server follows a simple version of those ideas.  Each method in a controller implements a request.  The methods take a standard set of arguments defined using interfaces rather than real objects.  The database repositories that supply data to the model also supply it as interfaces.  This allows objects that conform to the same interface to be used interchangeably during testing and thus supports TDD using dummy objects and mocks.  I've created objects modelled on the Struts form beans to carry the data to be rendered by the views.  Each form object can contain a notification message and/or an error message to be displayed at the top of the page, data items and error messages about the data items.  When the user submits data from an HTML form, it's handed to the controller method in a form object.

The requests and web pages are laid out using the REST model, implemented using the go-restful library.  A RESTful web server provides a set of resources that the user can access.  Each resource has its own model and controller, plus a set of views.  A resource can be (but need not be) represented by a database table.  All requests concerning a resource follow a pattern that starts with the resource name, for example:

```
    /people		display all people
    /people/42		display the person with ID 42
    /people/42/edit	fetch the data for person 42 and display a screen to change it
```

The server is called films because a future version will display information about films - a very simple form of IMDb.  The people table will hold data about actors, directors and so on, and there will be other tables, with web pages to manipulate them.  At present there is one resource backed by one table, so one model, one controller and one set of views.

The standard Go library includes a library net/html, which provides a framework for building and displaying web pages.  I use this to provide the views.  Each view takes a package of data provided by the controller, creates an HTML page to display it and sends the page to the user's browser.  The contents of the view is determined by a form object.

In Go, an object satisfies an interfaces if it implements all of the methods of the interface, so you can create an interface and then create an object that satisfies it, or you can take an existing object and write an interface that it satisfies.  Unlike with languages such as Java, this means that an interface can be retrofitted to an object that you did not create.  For example, the net/html package defines a structure called a Template which is used to render an HTML page.  The films server includes an interface Template that the HTML template structure satisfies.  The controller uses the templates in terms of that interface, so its templates can be replaced with mock versions during testing.

In the films server, all requests are sent to a function called marshal.  This figures out which controller to call (currently there's only one, the people controller).  It creates a services object which provides database repositories and the HTML templates, and creates a controller, setting the services object.  The controller then creates a suitable form filled in with data from the request and calls one of the controller's CRUD methods, passing the form as an interface.  The controller method uses the repositories to access data and the templates to render results.  So the controller methods are all driven by repositories, forms and templates supplied from outside.  The controller methods can be tested in isolation by supplying suitably-crafted dummy and mock objects.  I've written a number of example tests that do that.

The Database
============

The database runs under MySQL.  The server uses the GORP library to connect to the database.  That, of course, is the concern only of the model, and another model can be slotted in to replace it.  In the unit tests, the MySQL model is replaced by objects that provide fixed datasets that drive the logic through the desired logic paths.  Some of the tests use mocking provided by gomock and pegomock.

Testing
=======

The films server includes unit tests, each of which tests a single unit of software in isolation by providing it with dummy objects containing data specifically written for the test.  There are also integration tests, where a few units of software are bound together and tests are run to check that everything hangs together.  These tests check that the controller and the MySQL model work together correctly.  Finally, system tests check that the whole system works together.  Go includes a facility for running tests on a complete web server, but I use Selenium, specifically the Firefox Selenium addon, which can be used to test any web server, regardless of the technology used to create it.

The test directory contains a number of Selenium scripts for system testing.  To run them in Firefox, install the Selenium addon, start it up and use the file menu to load one of the test suites.  Use the green arrow buttons to run the whole suite or one of the tests in it.  As the test runs, you can see the results in your browser.

The unit and integration tests are written using the standard Go test facilities, plus gomock and pegomock for mocking.  They live in the same directory as the module that they are testing.  Go offers a naming convention that can be used to classify tests.  The test brace can be made to run only those tests whose names match a pattern.  I use this to differentiate between unit and integration tests.  On a UNIX or Linux server, the tests can be run using the shell script test.sh as follows:

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

At present I'm experimenting with using hand-crafted mocks and generating them automatically using the mocking frameworks gomock and pegomock, and I have tests that use all three techniques.  There are other mocking frameworks available.  The pegomock github page includes a survey of them.

It's worth saying that not everybody in the Go community agrees with the idea of using a mocking framework for test-driven development.  I find them a useful and cost-effective way to test my software.

Go also offers a framework that you can use to build a complete set of system tests for a web server, making web requests and checking that the correct response comes back.  You can use that, but of course, you still have to test your server using a web browser, preferably using all the common web browsers, to make sure that everything looks sensible.  I use Selenium for system testing web servers because it's visual, it can record tests interactively and it runs in a real browser.
