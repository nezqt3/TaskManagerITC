package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func CheckTelegramAuth(data map[string]string, botToken string) error {
	var checkList []string
	for k, v := range data {
		if k != "hash" {
			checkList = append(checkList, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(checkList)
	dataString := strings.Join(checkList, "\n")
	secretKey := sha256.Sum256([]byte(botToken))

	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if expectedHash != data["hash"] {
		return fmt.Errorf("invalid telegram auth hash")
	}
	return nil
}