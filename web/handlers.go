package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
	"web_service_GO/logger"
	"web_service_GO/pkg/task"
)

func (ds * DefaultServer)handleSubmit(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("New submit request")

	urlToUse := r.FormValue("url")

	byteID, err := exec.Command("uuidgen").Output() // use POSIX command to generate unique key
	if err != nil {
		logger.Error.Println("Can't generate uuid", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Can't generate id")
		return
	}
	uniqueID := string(byteID)
	req := task.UserRequest {
		ID:  uniqueID,
		URL: urlToUse,
	}

	ds.DB.Save(req)

	f, erBool := w.(http.Flusher)
	if erBool == false {
		logger.Error.Println("Flush error")
		fmt.Fprintln(w, "Flush error, can't print id")
		return
	}
	fmt.Fprintln(w, "Your id:", uniqueID)
	f.Flush()

	go ds.Calc.CalculateMD5(uniqueID, urlToUse) // each process starts in it's own goroutine
}

func (ds * DefaultServer)handleCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("New request to ", r.URL)
	id, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong request, should be /check/<id>")
	}
	request, err := ds.DB.Load(id)
	if err != nil {
		logger.Info.Println(id, " now found in database")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, id, "not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, request)
}

func (ds * DefaultServer)handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Wrong path, use /check /submit instead")
}