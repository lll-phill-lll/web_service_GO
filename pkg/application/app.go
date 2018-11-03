package application

import (
	"web_service_GO/pkg/DB"
	"web_service_GO/pkg/calc"
	"web_service_GO/web"
)

type App struct {
	DB        DB.Database
	Server    web.Server
	Processor calc.Calc
}

func (app *App) Start() {
	app.Server.SetEndpoints()
	app.Server.StartServe(8080)
}
