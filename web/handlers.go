package web

import (
	"encoding/json"
	"fmt"
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

	go ds.Calc.CalculateMD5(uniqueID) // each process starts in it's own goroutine
}

func (ds * DefaultServer)handleCheck(w http.ResponseWriter, r *http.Request) {
	myParam := r.URL.Query().Get("id")

	if myParam != "" {
		myParam += "\n"
		mu.Lock()
		i, inMap := allRequests[myParam]
		mu.Unlock()

		if inMap {
			if i.ready {
				if i.er {
					fmt.Fprintln(w, "Error during md5 computing")
				} else {
					fmt.Fprintln(w, "{md5:", i.md5, ", status: done, url:", i.url, "}")
				}
			} else {
				fmt.Fprintln(w, "{status : running}")
			}
		} else {
			fmt.Fprintln(w, "{Not exist}")
		}
	} else {
		fmt.Fprintf(w, "incorrect request. Id error.")
	}
}

func (ds * DefaultServer)handleRoot(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t api.UserRequest
	if err := decoder.Decode(&t); err != nil {
		logger.Error.Println("Status:", http.StatusBadRequest, err)
		sendResponse(http.StatusBadRequest, "Bad Request", r, w)
		return
	}

	logger.Info.Println("Status:", http.StatusAccepted, "Request:", r.Host, r.URL.Path)
	sendResponse(http.StatusAccepted, "Accepted", r, w)

	go docker.Run(t, r.Host, r.URL.Path) // use go to run process ib new goroutine
}