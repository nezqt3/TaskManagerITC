package model

type ProjectMember struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type Project struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Members     []ProjectMember `json:"members"`
}
