package main

import (
	"encoding/hex"
	"fmt"
	"crypto/md5"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time" // to get "running"" status
	"web_service_GO/logger"
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
	db := &DB.Database{}
	server := &serv.DefaultServer{}
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
