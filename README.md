# films
A simple MVC web server written in Go with restful interfaces providing CRUD operations on database tables

This server is deliberately very simple and is for demonstration prposes.  It currently supports a MySQL database called films with one table people.

The database should be accessible by a user called webuser with password "secret".  (This is defined in MakeGorpMysqlDBSession() in utilities/dbsession/dbsession.go.)

Start the server from a command window as follows:

In the top-level diretory, run setanv.sh like so:

    .  setenv.sh

(The "." at the start is a command and must be present.)

Change to the directory containing the views directory.

Run the server like so

     films

The server listems on port 4000.  In a web browser, navigate to

    http://localhost:4000/people

Initially there will be no entries in the table.  Use the create button to create some.

Each entry has a forename and a surname.  In the create screen, try missing one or both of them out and pressing the submit button.
