package main

import (
	"os"
	"web_service_GO/logger"
	"web_service_GO/pkg/DB"
	"web_service_GO/pkg/application"
	"web_service_GO/serv"
)

type userRequest struct { // unique struct for each request
	md5   string
	url   string
	ready bool
	id    string
	er    bool
}


func InitApp() application.App {
	db := &DB.MapDatabase{}
	server := &serv.DefaultServer{
		DB: db,
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
