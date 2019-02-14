package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/jcorry/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()

	td.Flash = app.session.PopString(r, "flash")

	td.AuthenticatedUser = app.authenticatedUser(r)
	return td
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
