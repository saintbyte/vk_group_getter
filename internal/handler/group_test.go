package handler

import "testing"

func TestNormalizeAddress(t *testing.T) {
	tests := []struct {
		input    string
		wantAddr string
		wantType string
	}{
		{"https://vk.com/durov", "durov", "vk_group"},
		{"http://vk.com/durov", "durov", "vk_group"},
		{"vk.com/durov", "durov", "vk_group"},
		{"m.vk.com/durov", "durov", "vk_group"},
		{"https://m.vk.com/durov", "durov", "vk_group"},
		{"http://m.vk.com/durov", "durov", "vk_group"},
		{"https://vk.ru/durov", "durov", "vk_group"},
		{"http://vk.ru/durov", "durov", "vk_group"},
		{"vk.ru/durov", "durov", "vk_group"},
		{"m.vk.ru/durov", "durov", "vk_group"},
		{"https://m.vk.ru/durov", "durov", "vk_group"},
		{"http://m.vk.ru/durov", "durov", "vk_group"},
		{"https://vk.com/durov/", "durov", "vk_group"},
		{"https://vk.com/durov   ", "durov", "vk_group"},
		{"https://VK.COM/Durov", "Durov", "vk_group"},
		{"durov", "durov", "vk_group"},
		{"", "", "vk_group"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotAddr, gotType := normalizeAddress(tt.input)
			if gotAddr != tt.wantAddr {
				t.Errorf("normalizeAddress(%q) addr = %q, want %q", tt.input, gotAddr, tt.wantAddr)
			}
			if gotType != tt.wantType {
				t.Errorf("normalizeAddress(%q) type = %q, want %q", tt.input, gotType, tt.wantType)
			}
		})
	}
}
