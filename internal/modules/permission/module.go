package permission

import "letsgo/internal/app"

type Module struct {
	handler *Handler
}

func NewModule(application *app.Application) *Module {
	return &Module{
		handler: NewHandler(NewService(NewRepository(application))),
	}
}
