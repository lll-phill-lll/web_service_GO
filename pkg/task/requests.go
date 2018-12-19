package task

type UserRequest struct {
	MD5    string `json:"md5"`
	URL    string `json:"url"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

var RequestStatus = struct {
	Waiting string
	Failed  string
	Ready   string
	Running string
}{
	"waiting",
	"failed",
	"ready",
	"running",
}
