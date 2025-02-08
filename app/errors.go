package app

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	app.ClientError(w, http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	ts, ok := app.templateCache["error.html"]
	if !ok {
		err := fmt.Errorf("the template \"error\" does not exist")
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data := struct {
		ErrorCode int
		ErrorText string
	}{
		ErrorCode: status,
		ErrorText: http.StatusText(status),
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "errorBase", data)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	return
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
