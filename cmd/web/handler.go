package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// http.NotFound(w, r) --> replacing this with helper function to achieve centralized error handling
		app.notFound(w) //notFound() helper
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err) // using the serverError() helper function
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err) //serverError()
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w) // notFound() helper function
		return
	}

	fmt.Fprintf(w, "Displaying the specific snippet with ID: %d", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allowed", "POST")
		// http.Error(w, "Method Not Allowed", 405)
		app.clientError(w, http.StatusNotFound)
		return
	}

	w.Write([]byte("Create a new Snippet"))
}
