package model

type UserCSV struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	DateOfBirthday string `json:"date_of_birthday"`
	NumberOfPhone string `json:"number_of_phone"`
	Role string `json:"role"`
	TelegramID string `json:"telegram_id"`
	MayToOpen bool `json:"may_to_open"`
}