package filemanager

type (
	Clients struct {
		Uuid     string `json:"uuid"`
		Username string `json:"username"`
		Qoute    string `json:"qoute"`
	}
	User struct {
		UUID     string `json:"uuid"`
		UserName string `json:"username"`
	}
)
