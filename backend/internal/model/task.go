package model

type Task struct {
	ID                int    `json:"id"`
	Description       string `json:"description"`
	Deadline          string `json:"deadline"`
	Status            string `json:"status"`
	CompletionMessage string `json:"completion_message,omitempty"`
	ReviewMessage     string `json:"review_message,omitempty"`
	ReviewedBy        string `json:"reviewed_by,omitempty"`
	ReviewedAt        string `json:"reviewed_at,omitempty"`
	User              string `json:"user"`
	Title             string `json:"title"`
	Author            string `json:"author"`
	IdProject         int    `json:"id_project"`
	IdUser            int64  `json:"id_user"`
	ProjectTitle      string `json:"project_title,omitempty"`
}
