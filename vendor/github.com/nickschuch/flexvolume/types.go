package flexvolume

type Status string

const (
	StatusSuccess Status = "Success"
	StatusFailure Status = "Failure"
)

type FlexVolume interface {
	Init() Response
	Attach(map[string]string) Response
	Detach(string) Response
	Mount(string, string, map[string]string) Response
	Unmount(string) Response
}

type Response struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
	Device  string `json:"device"`
}
