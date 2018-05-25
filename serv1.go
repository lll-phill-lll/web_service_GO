package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sync"
	// "time" // to get "running"" status
)

type userRequest struct { // unique struct for each request
	md5   string
	url   string
	ready bool
	id    string
	er    bool
}

var mu = &sync.Mutex{} // add mutex to avoid race condition (to check use -race flag while compiling)

var allRequests = make(map[string]userRequest) // map to store all the requests

func handlerCheck(w http.ResponseWriter, r *http.Request) {
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

func startMD5(url string, uniqueID string) {
	response, er := http.Get(url)
	if er != nil {
		fmt.Println("Get url error. ID=", uniqueID)
	}

	var body []byte
	if er == nil {
		defer response.Body.Close()

		body, er = ioutil.ReadAll(response.Body)
		if er != nil {
			fmt.Println("Error while getting file body. ID=", uniqueID)
		}
	}

	hasher := md5.New()
	if er == nil {
		hasher.Write(body)
	}

	// time.Sleep(25 * time.Second) // to get "running" status
	mu.Lock()
	thisRequest := allRequests[uniqueID]
	thisRequest.ready = true
	if er == nil { // check if there errors while saving file and computing md5
		thisRequest.md5 = hex.EncodeToString(hasher.Sum(nil))
	} else {
		thisRequest.er = true
	}
	delete(allRequests, uniqueID)
	allRequests[uniqueID] = thisRequest
	mu.Unlock()
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "Should be method POST")
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
	if erBool == true {
		fmt.Fprintln(w, "Flush error, can't print id")
		return
	}
	f.Flush()

	go startMD5(urlToUse, uniqueID) // each process starts in it's own goroutine

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "use /submit or /check prefix")
}

func main() {
	http.HandleFunc("/submit", handleSubmit)
	http.HandleFunc("/check", handlerCheck)
	http.HandleFunc("/", handleRoot)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)

}
