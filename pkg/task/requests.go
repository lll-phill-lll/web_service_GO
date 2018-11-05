package task

type UserRequest struct {
	MD5   string	`json:"md5"`
	URL   string	`json:"url"`
	ID    string	`json:"id"`
	Status string	`json:"status"`
}

var RequestStatus = struct {
	Failed string
	Ready string
	Running string

}{
	"failed",
	"ready",
	"running",
}
