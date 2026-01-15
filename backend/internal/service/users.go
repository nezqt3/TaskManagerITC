package service

import (
	"fmt"
	"net/http"
	"encoding/csv"
	"log"

	"backend/internal/model"
)

func GetUsers(cfg *model.Config) ([]model.UserCSV, error) {
	resp, err := http.Get(cfg.SPREADSHEET_URL)
	var users []model.UserCSV

	if err != nil {
		log.Fatalf("Ошибка при получении пользователей")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode = %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV пустой")
	}

	for i, row := range records {
		if i == 0 || i == 1 {
			continue
		}

		mayToOpen := row[6] == "TRUE"

		user := model.UserCSV{
			FullName: row[0],
			Username: row[1],
			DateOfBirthday: row[2],
			NumberOfPhone: row[3],
			Role: row[4],
			TelegramID: row[5],
			MayToOpen: mayToOpen,
		}

		users = append(users, user)
	}

	return users, err
}