package media

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Valid FTP URL", "ftp://example.com", true},
		{"Valid URL with path", "https://example.com/path/to/resource", true},
		{"Valid URL with query parameters", "https://example.com?query=param", true},
		{"Valid URL with port", "http://example.com:8080", true},
		{"Invalid URL - missing scheme", "example.com", false},
		{"Invalid URL - missing host", "http://", false},
		{"Invalid URL - malformed", "://example.com", false},
		{"Invalid URL - empty string", "", false},
		//{"Invalid URL - spaces", "http://example.com/with spaces", false},
		//{"Invalid URL - invalid characters", "http://example.com/with^invalid^chars", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidURL(tt.link)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v; expected %v", tt.link, result, tt.expected)
			}
		})
	}
}
