package domain

type User struct {
	Uuid    string `json:"uuid"`
	Ip      string `json:"ip"`
	Os      string `json:"os"`
	Browser string `json:"browser"`
}
