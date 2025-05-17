package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Logs Errors + Stack trace, shows user 500 error
func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends Client errors like 400, 405 etc
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Wrapper for 404 Page not found
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
