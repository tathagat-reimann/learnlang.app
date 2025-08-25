package utils

import "testing"

func TestMakePackKey(t *testing.T) {
	tests := []struct {
		userID   string
		langID   string
		packName string
		expected string
	}{
		{"User1", "EN", "Starter", "user1:en:starter"},
		{"ADMIN", "De", "Pro", "admin:de:pro"},
		{"", "", "", ""},
	}

	for _, tt := range tests {
		result := MakePackKey(tt.userID, tt.langID, tt.packName)
		if result != tt.expected {
			t.Errorf("MakePackKey(%q, %q, %q) = %q; want %q", tt.userID, tt.langID, tt.packName, result, tt.expected)
		}
	}
}
