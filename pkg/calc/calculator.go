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
	CalculateMD5(string)
}

type DefaultCalc struct {
	DB DB.Database
}

func (dc *DefaultCalc) CalculateMD5(id string) {
	currentTask, err := dc.DB.Load(id)
	if err != nil {
		logger.Error.Println("Task not found in database, id:", id)
		currentTask = task.UserRequest{
			ID:     id,
			URL:    "",
			Status: task.RequestStatus.Failed,
		}
		return
	}
	currentTask.Status = task.RequestStatus.Running
	dc.DB.Save(currentTask)
	response, err := http.Get(currentTask.URL)
	if err != nil {
		logger.Error.Println("Can't load file. Error:", err, "ID =", id, "URL = ", currentTask.URL)
		currentTask.Status = task.RequestStatus.Failed
		dc.DB.Save(currentTask)
		return
	}

	var body []byte
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error.Println("Error while getting file body. ID=", id)
		currentTask.Status = task.RequestStatus.Failed
		dc.DB.Save(currentTask)
		return
	}

	hasher := md5.New()
	hasher.Write(body)

	currentTask.MD5 = hex.EncodeToString(hasher.Sum(nil))
	currentTask.Status = task.RequestStatus.Ready
	logger.Info.Println("MD5 computed for id=", currentTask.ID)
	dc.DB.Save(currentTask)
}
