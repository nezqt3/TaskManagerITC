package model

type Event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	TimeRange   string `json:"time_range"`
	CreatedBy   string `json:"created_by"`
	Description string `json:"description,omitempty"`
}
