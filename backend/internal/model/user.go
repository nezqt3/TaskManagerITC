package model

type UserProfile struct {
	TelegramID     string `json:"telegram_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name,omitempty"`
	Username       string `json:"username,omitempty"`
	PhotoURL       string `json:"photo_url,omitempty"`
	FullName       string `json:"full_name"`
	DateOfBirthday string `json:"date_of_birthday"`
	NumberOfPhone  string `json:"number_of_phone"`
	Role           string `json:"role"`
	MayToOpen      bool   `json:"may_to_open"`
}
