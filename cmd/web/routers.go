package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(app.recoverPanic)
	router.Use(app.logRequest)
	router.Use(secureHeaders)
	router.Use(app.session.LoadAndSave)

	//Public routes
	router.Get("/", app.home)
	router.Get("/snippet/{id}", app.showSnippet)

	//Auth-only Routes
	router.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)
		r.Get("/snippet/create", app.createSnippetForm)
		r.Post("/snippet/create", app.createSnippet)
		r.Post("/user/logout", app.logoutUser)
	})

	//Auth routes
	router.Get("/user/signup", app.signupUserForm)
	router.Post("/user/signup", app.signupUser)
	router.Get("/user/login", app.loginUserForm)
	router.Post("/user/login", app.loginUser)

	//Static files
	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	return router
}
