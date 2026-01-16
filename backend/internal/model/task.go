package model

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
	User        string `json:"user"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	IdProject   int    `json:"id_project"`
	IdUser      int64  `json:"id_user"`
}