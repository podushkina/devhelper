package encoder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_Base64Encode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		urlSafe     bool
		expected    string
		expectError bool
	}{
		{
			name:        "Basic encode",
			input:       "Hello, World!",
			urlSafe:     false,
			expected:    "SGVsbG8sIFdvcmxkIQ==",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			urlSafe:     false,
			expected:    "",
			expectError: false,
		},
		{
			name:        "URL-safe encode",
			input:       "Hello, World!+/",
			urlSafe:     true,
			expected:    "SGVsbG8sIFdvcmxkISsv",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			encoder := NewEncoder(in, out)

			err := encoder.Base64Encode(tt.urlSafe)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Удаляем перевод строки для сравнения
				result := strings.TrimSpace(out.String())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestEncoder_Base64Decode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		urlSafe     bool
		expected    string
		expectError bool
	}{
		{
			name:        "Basic decode",
			input:       "SGVsbG8sIFdvcmxkIQ==",
			urlSafe:     false,
			expected:    "Hello, World!",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			urlSafe:     false,
			expected:    "",
			expectError: false,
		},
		{
			name:        "URL-safe decode",
			input:       "SGVsbG8sIFdvcmxkISsv",
			urlSafe:     true,
			expected:    "Hello, World!+/",
			expectError: false,
		},
		{
			name:        "Invalid Base64",
			input:       "Invalid Base64!@#",
			urlSafe:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "With whitespace",
			input:       "  SGVsbG8sIFdvcmxkIQ==  ",
			urlSafe:     false,
			expected:    "Hello, World!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			encoder := NewEncoder(in, out)

			err := encoder.Base64Decode(tt.urlSafe)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, out.String())
			}
		})
	}
}

func TestEncoder_URLEncode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Basic URL encode",
			input:       "Hello, World!",
			expected:    "Hello%2C+World%21",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Special characters",
			input:       "test?query=value&param=data",
			expected:    "test%3Fquery%3Dvalue%26param%3Ddata",
			expectError: false,
		},
		{
			name:        "UTF-8 characters",
			input:       "тест",
			expected:    "%D1%82%D0%B5%D1%81%D1%82",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			encoder := NewEncoder(in, out)

			err := encoder.URLEncode()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Удаляем перевод строки для сравнения
				result := strings.TrimSpace(out.String())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestEncoder_URLDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Basic URL decode",
			input:       "Hello%2C+World%21",
			expected:    "Hello, World!",
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Special characters",
			input:       "test%3Fquery%3Dvalue%26param%3Ddata",
			expected:    "test?query=value&param=data",
			expectError: false,
		},
		{
			name:        "UTF-8 characters",
			input:       "%D1%82%D0%B5%D1%81%D1%82",
			expected:    "тест",
			expectError: false,
		},
		{
			name:        "Invalid URL encoding",
			input:       "%ZZ",
			expected:    "",
			expectError: true,
		},
		{
			name:        "With whitespace",
			input:       "  Hello%2C+World%21  ",
			expected:    "Hello, World!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			encoder := NewEncoder(in, out)

			err := encoder.URLDecode()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Удаляем перевод строки для сравнения
				result := strings.TrimSpace(out.String())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNewEncoder(t *testing.T) {
	in := strings.NewReader("test")
	out := new(bytes.Buffer)

	encoder := NewEncoder(in, out)

	assert.NotNil(t, encoder)
	assert.Equal(t, in, encoder.reader)
	assert.Equal(t, out, encoder.writer)
}
