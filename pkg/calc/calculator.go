package calc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"web_service_GO/pkg/DB"
)

type Calc interface {
	CalculateMD5(string)
}

type DefaultCalc struct {
	db * DB.Database
}

func (dc * DefaultCalc) CalculateMD5(url string) {
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

	time.Sleep(0 * time.Second) // to get "running" status, change 0 to 25 while testing
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
