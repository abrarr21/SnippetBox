package main

import (
	"bytes"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error, call our serverError helper and return
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
	}

	// Write the content of the buffer to the http.ResponseWriter. Again, this is another time where we pass our http.ResponseWriter to a function that takes an io.Writer
	buf.WriteTo(w)
}
