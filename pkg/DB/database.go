package DB

import (
	"errors"
	"sync"
	"web_service_GO/pkg/task"
)

type Database interface {
	Load(string) (task.UserRequest, error)
	Save(task.UserRequest)
}

type MapDatabase struct {
	sync.RWMutex
	db map[string]task.UserRequest
}

func NewMapDataBase() *MapDatabase {
	return &MapDatabase{db: map[string]task.UserRequest{}}
}

func (b *MapDatabase) Load(id string) (task.UserRequest, error) {
	b.RLock()
	defer b.RUnlock()
	request, in := b.db[id]
	if in {
		return request, nil
	}
	return request, errors.New("not exist")
}

func (b *MapDatabase) Save(request task.UserRequest) {
	id := request.ID
	b.Lock()
	defer b.Unlock()
	b.db[id] = request
}
