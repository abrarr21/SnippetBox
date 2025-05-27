package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abrarr21/snippet/pkg/models"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		// r.RemoteAddr --> client ka IP + PORT
		// r.Proto --> HTTP version (like http/1.1)
		// r.Method --> GET,POST etc
		// r.URL.RequestURI --> requested paths + queries
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if a userID value exists in the session. If this isn't present then call the next handler in the chain as normal.
		exists := app.session.Exists(r.Context(), "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		//Fetch the details of the current user from the Database. If no matching record is found, remove the (invalid) userID from their session and call the next handler in the chain as normal
		user, err := app.users.Get(app.session.GetInt(r.Context(), "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r.Context(), "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		//Otherwise, we know that the request is coming from the valid, authenticated user. we create a new copy of the request with the user information added to the request context, and call the next handler  in the chian *using this new copy of the request.
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
