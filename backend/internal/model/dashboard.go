package model

type Dashboard struct {
	Projects []Project `json:"projects"`
	Tasks    []Task    `json:"tasks"`
	Events   []Event   `json:"events"`
}
