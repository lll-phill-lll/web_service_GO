package calc

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"web_service_GO/logger"
	"web_service_GO/pkg/DB"
	"web_service_GO/pkg/task"
)

type Calc interface {
	CalculateMD5(string, string)
}

type DefaultCalc struct {
	DB DB.Database
}

func (dc * DefaultCalc) CalculateMD5(id string, url string) {
	request := task.UserRequest{
		ID: id,
		URL: url,
		Status: task.RequestStatus.Running,
	}
	dc.DB.Save(request)
	response, err := http.Get(url)
	if err != nil {
		logger.Error.Println("Get url error. ID=", id)
		request.Status = task.RequestStatus.Failed
		dc.DB.Save(request)
		return
	}


	var body []byte
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error.Println("Error while getting file body. ID=", id)
		request.Status = task.RequestStatus.Failed
		dc.DB.Save(request)
		return
	}

	hasher := md5.New()
	hasher.Write(body)

	request.MD5 = hex.EncodeToString(hasher.Sum(nil))
	request.Status = task.RequestStatus.Ready
	dc.DB.Save(request)
}
