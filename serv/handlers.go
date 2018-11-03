package serv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

func (ds * DefaultServer)handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "Should be method POST with /submit")
		return
	}

	urlToUse := r.FormValue("url")

	byteID, er := exec.Command("uuidgen").Output() // use POSIX command to generate unique key
	if er != nil {
		fmt.Fprintln(w, "uuidgen error, can't generate id")
		return
	}
	uniqueID := string(byteID)
	var req userRequest = userRequest{
		id:  uniqueID,
		url: urlToUse,
	}

	mu.Lock()
	allRequests[uniqueID] = req
	mu.Unlock()

	f, erBool := w.(http.Flusher)
	if erBool == false {
		fmt.Fprintln(w, "Flush error, can't print id")
		return
	}
	fmt.Fprintln(w, "your id:", uniqueID)
	f.Flush()

	go startMD5(urlToUse, uniqueID) // each process starts in it's own goroutine
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