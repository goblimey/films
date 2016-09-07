package utilities

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/goblimey/films/retrofit/restful/response"
)

// Dead displays a hand-crafted error page.  It's the page of last resort.
func Dead(resp response.Response) {
	log.SetPrefix("Dead() ")
	log.Println()
	defer noPanic()
	fmt.Sprintf("foo", "1", "2")
	html := fmt.Sprintf("%s%s%s%s%s%s\n",
		"<html><head></head><body>",
		"<p><b><font color=\"red\">",
		"This server is experiencing a Total Inability To Service Usual Processing (TITSUP).",
		"</font></b></p>",
		"<p>We will be restoring normality just as soon as we are sure what is normal anyway.</p>",
		"</body></html>")

	_, err := fmt.Fprintln(resp.Response().ResponseWriter, html)
	if err != nil {
		log.Printf("error while attempting to display the error page of last resort - %s", err.Error())
		http.Error(resp.Response().ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

// Recover from any panic and log an error.
func noPanic() {
	if p := recover(); p != nil {
		log.Printf("unrecoverable internal error %v\n", p)
	}
}

// Trim removes leading and trailing white space from a string.
func Trim(str string) string {
	return strings.Trim(str, " \t\n")
}

// Map2String displays the contents of a map of strings with string values as a
// single string.The field named "foo" with value "bar" becomes 'foo="bar",'.
func Map2String(m map[string]string) string {
	// The result array has two entries for each map key plus leading and
	// trailing brackets.
	result := make([]string, 0, 2+len(m)*2)
	result = append(result, "[")
	for key, value := range m {
		result = append(result, key)
		result = append(result, "=\"")
		result = append(result, value)
		result = append(result, "\",")
	}
	result = append(result, "]")

	return strings.Join(result, "")
}
