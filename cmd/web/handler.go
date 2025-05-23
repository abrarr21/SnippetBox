package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/abrarr21/snippet/pkg/models"
	"github.com/go-chi/chi"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r) --> replacing this with helper function to achieve centralized error handling
	// 	app.notFound(w) //notFound() helper
	// 	return
	// }

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// data := &templateData{Snippets: s}
	//
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	//
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
	//
	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	//Use the new render helper to replace above commented code
	app.render(w, r, "home.page.html", &templateData{Snippets: s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w) // notFound() helper function
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// data := &templateData{Snippet: s}
	//
	// files := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	app.serverError(w, err) // using the serverError() helper function
	// 	return
	// }
	//
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	app.serverError(w, err) //serverError()
	// }
	//
	// fmt.Fprintf(w, "%v", s)

	// Use the render helper to replace above commented code
	app.render(w, r, "show.page.html", &templateData{Snippet: s})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	// http.Error(w, "Method Not Allowed", 405)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	title := "This is the title"
	content := "content of the snippet"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", nil)
}
