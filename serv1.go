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

type userRequest struct {
	md5   string
	url   string
	ready bool
	id    string
	er    bool
}

var mu = &sync.Mutex{} // add mutex to avoid race condition (to check use -race flag while compiling)

var allRequests = make(map[string]userRequest)

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

func startMD5(url string, uniqueID string) error {
	response, _ := http.Get(url)

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	hasher := md5.New()
	hasher.Write(body)

	// time.Sleep(25 * time.Second) // to get "running" status
	mu.Lock()
	thisRequest := allRequests[uniqueID]
	thisRequest.ready = true
	thisRequest.md5 = hex.EncodeToString(hasher.Sum(nil))
	delete(allRequests, uniqueID)
	allRequests[uniqueID] = thisRequest
	mu.Unlock()
	return nil
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "Should be method POST")
		return
	}

	urlToUse := r.FormValue("url")

	byteID, _ := exec.Command("uuidgen").Output()
	uniqueID := string(byteID)
	var req userRequest = userRequest{
		id:  uniqueID,
		url: urlToUse,
	}

	mu.Lock()
	allRequests[uniqueID] = req
	mu.Unlock()

	fmt.Fprintln(w, "your id:", uniqueID)

	f, _ := w.(http.Flusher)
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
