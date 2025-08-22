package utils

import (
	"fmt"
	"strings"
)

func MakePackKey(userID, langCode, packName string) string {
	return fmt.Sprintf("%s:%s:%s", userID, langCode, strings.ToLower(packName))
}
