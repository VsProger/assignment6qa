package handlers

import (
	"forum/app"
	"forum/internal/service"
)

type handler struct {
	service service.ServiceI
	app     *app.Application
}

func New(s service.ServiceI, app *app.Application) *handler {
	return &handler{
		s,
		app,
	}
}
