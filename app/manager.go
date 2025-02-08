package app

import (
	"html/template"
	"log"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	templateCache map[string]*template.Template
	// snippets       models.SnippetModelInterface
	// users          models.UserModelInterface

	// formDecoder    *form.Decoder
	// sessionManager *scs.SessionManager
}

func New(infoLog, errorLog *log.Logger, templateCache map[string]*template.Template) *Application {
	return &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		templateCache: templateCache,
	}
}
