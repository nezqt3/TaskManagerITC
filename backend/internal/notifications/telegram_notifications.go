package notifications

import (
	"net/http"
	"net/url"
	"io"

	"fmt"
	"backend/internal/model"
	"backend/internal/logger"
)

func SendTelegramNotification(cfg *model.Config, id int64, text string) error{
	msg := url.QueryEscape(text)
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&parse_mode=HTML", cfg.TelegramBotToken, id, msg)

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)

	logger.Info.Printf(
		"Telegram response: status=%d body=%s",
		resp.StatusCode,
		string(body),
	)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: %s", resp.Status)
	}

	return nil
}