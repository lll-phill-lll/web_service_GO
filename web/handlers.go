package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os/exec"
	"strings"
	"web_service_GO/logger"
	"web_service_GO/pkg/task"
)

func (ds *DefaultServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	urlToUse := r.FormValue("url")

	byteID, err := exec.Command("uuidgen").Output() // use POSIX command to generate unique key
	if err != nil {
		logger.Error.Println("Can't generate uuid, to request:", urlToUse, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Can't generate id")
		return
	}
	uniqueID := strings.TrimSpace(string(byteID))

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Your id:", uniqueID)

	f, erBool := w.(http.Flusher)
	if erBool == false {
		logger.Error.Println("Flush error")
		fmt.Fprintln(w, "Flush error, can't print id")
		return
	}
	f.Flush()

	logger.Info.Println("New submit request, id = ", uniqueID)

	req := task.UserRequest{
		ID:  uniqueID,
		URL: urlToUse,
	}

	ds.DB.Save(req)

	go ds.Calc.CalculateMD5(uniqueID, urlToUse) // each process starts in it's own goroutine
}

func (ds *DefaultServer) handleCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("New request to ", r.URL)
	id, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong request, should be /check/<id>")
	}
	request, err := ds.DB.Load(id)
	if err != nil {
		logger.Info.Println(id, "not found in database")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, id, "not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(request)
	if err != nil {
		logger.Error.Println("Can't marshal json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(ret))
}

func (ds *DefaultServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Wrong path, use /check /submit instead")
}
