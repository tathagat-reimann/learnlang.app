package utils

import (
	"fmt"
	"strings"
)

func MakePackKey(userID, langID, packName string) string {
	if userID == "" || langID == "" || packName == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s", strings.ToLower(userID), strings.ToLower(langID), strings.ToLower(packName))
}

// MakeVocabKey creates a composite key to uniquely identify a vocab entry
// for a given user, language, pack and word text (case-insensitive).
func MakeVocabKeyByPackID(packID, name string) string {
	if packID == "" || name == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", strings.ToLower(packID), strings.ToLower(name))
}
