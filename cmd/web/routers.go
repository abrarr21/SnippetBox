package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// using "justinas/alice package for better middleware chaining"
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	r := chi.NewRouter()

	r.Get("/", app.home)
	r.Get("/snippet/{id}", app.showSnippet)
	r.Post("/snippet/create", app.createSnippet)
	r.Get("/snippet/create", app.createSnippetForm)

	// Create a file server which serves the files out of the "./ui/static/" directory.
	// Note that the path given to the http.Dir function is relative to the project's directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	// all URL paths that start with /static/. For matching paths, we strip the "/static" prefix before the request reaches the server.
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standardMiddleware.Then(r)
}
