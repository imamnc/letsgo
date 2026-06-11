package health

import "letsgo/internal/app"

type Module struct {
	handler *Handler
}

func NewModule(app *app.Application) *Module {
	return &Module{
		handler: NewHandler(NewService(NewRepository(app))),
	}
}
