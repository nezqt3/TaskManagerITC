package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"backend/internal/logger"
)

func CheckTelegramAuth(data map[string]string, botToken string) error {
	logger.Info.Println("Checking Telegram auth...")

	var checkList []string
	for k, v := range data {
		if k != "hash" {
			checkList = append(checkList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(checkList)
	dataString := strings.Join(checkList, "\n")
	logger.Info.Printf("Data string for hash: %s\n", dataString)

	secretKey := sha256.Sum256([]byte(botToken))

	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if expectedHash != data["hash"] {
		logger.Error.Printf("Telegram auth failed: invalid hash. Expected %s, got %s\n", expectedHash, data["hash"])
		return fmt.Errorf("invalid telegram auth hash")
	}

	logger.Info.Println("Telegram auth hash valid")
	return nil
}