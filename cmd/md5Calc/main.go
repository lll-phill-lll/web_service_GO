package main

import (
	"os"
	"web_service_GO/logger"
	"web_service_GO/pkg/DB"
	"web_service_GO/pkg/application"
	"web_service_GO/pkg/calc"
	"web_service_GO/web"
)



func InitApp() application.App {
	db := &DB.MapDatabase{}
	md5Calc := &calc.DefaultCalc{}
	server := &web.DefaultServer{
		DB: db,
		Calc: md5Calc,
	}
	app := application.App{
		DB: db,
		Server: server,
	}
	return app
}

func main() {
	// choose streams for each type of logs: stderr, stdout
	logger.SetLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	app := InitApp()
	app.Start()
}
