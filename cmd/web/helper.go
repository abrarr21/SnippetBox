package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/abrarr21/snippet/pkg/models"
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

// Create an addDefaultData helper. This takes a pointer to a templateData struct, adds the current year to the CurrentYear field, and then returns the pointer. Again, we're not using the *http.Request parameter at the moment, but we will do later in the book.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.popFlash(r) // Used popFlash helper function here.
	return td
}

// this function is creaetd for removing flash message
func (app *application) popFlash(r *http.Request) string {
	flash := app.session.GetString(r.Context(), "flash")
	if flash != "" {
		app.session.Remove(r.Context(), "flash")
	}

	return flash
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
	err := ts.Execute(buf, app.addDefaultData(td, r)) // Execute the template set, passing the dynamic data with current year injected.
	if err != nil {
		app.serverError(w, err)
	}

	// Write the content of the buffer to the http.ResponseWriter. Again, this is another time where we pass our http.ResponseWriter to a function that takes an io.Writer
	buf.WriteTo(w)
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
