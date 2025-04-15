package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		indent      int
		color       bool
		expectError bool
	}{
		{
			name:        "Valid JSON",
			input:       `{"name":"John","age":30,"city":"New York"}`,
			indent:      2,
			color:       false,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       `{"name":"John","age":30,"city":"New York"`,
			indent:      2,
			color:       false,
			expectError: true,
		},
		{
			name:        "Empty JSON",
			input:       `{}`,
			indent:      2,
			color:       false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			formatter := NewFormatter(in, out)

			err := formatter.FormatJSON(tt.indent, tt.color)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, out.String())
			}
		})
	}
}

func TestFormatYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		color       bool
		expectError bool
	}{
		{
			name: "Valid YAML",
			input: `
name: John
age: 30
city: New York
`,
			color:       false,
			expectError: false,
		},
		{
			name: "Invalid YAML",
			input: `
name: John
  age: 30
 city: New York
`,
			color:       false,
			expectError: true,
		},
		{
			name:        "Empty YAML",
			input:       `{}`,
			color:       false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			formatter := NewFormatter(in, out)

			err := formatter.FormatYAML(tt.color)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, out.String())
			}
		})
	}
}

func TestFormatXML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		indent      int
		color       bool
		expectError bool
	}{
		{
			name:        "Valid XML",
			input:       `<person><name>John</name><age>30</age><city>New York</city></person>`,
			indent:      2,
			color:       false,
			expectError: false,
		},
		{
			name:        "Invalid XML",
			input:       `<person><name>John</name><age>30</age><city>New York</city>`,
			indent:      2,
			color:       false,
			expectError: true,
		},
		{
			name:        "Empty XML",
			input:       `<root></root>`,
			indent:      2,
			color:       false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			formatter := NewFormatter(in, out)

			err := formatter.FormatXML(tt.indent, tt.color)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, out.String())
			}
		})
	}
}
