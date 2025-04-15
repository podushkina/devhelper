package converter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		input       string
		indent      int
		expectError bool
	}{
		{
			name:        "JSON to YAML - Valid",
			from:        "json",
			to:          "yaml",
			input:       `{"name":"John","age":30,"city":"New York"}`,
			indent:      2,
			expectError: false,
		},
		{
			name:        "JSON to YAML - Invalid JSON",
			from:        "json",
			to:          "yaml",
			input:       `{"name":"John","age":30,"city":"New York"`,
			indent:      2,
			expectError: true,
		},
		{
			name:        "YAML to JSON - Valid",
			from:        "yaml",
			to:          "json",
			input:       "name: John\nage: 30\ncity: New York",
			indent:      2,
			expectError: false,
		},
		{
			name:        "YAML to JSON - Invalid YAML",
			from:        "yaml",
			to:          "json",
			input:       "name: John\n  age: 30\n city: New York",
			indent:      2,
			expectError: true,
		},
		{
			name:        "XML to JSON - Valid",
			from:        "xml",
			to:          "json",
			input:       `<person><name>John</name><age>30</age><city>New York</city></person>`,
			indent:      2,
			expectError: false,
		},
		{
			name:        "XML to JSON - Invalid XML",
			from:        "xml",
			to:          "json",
			input:       `<person><name>John</name><age>30</age><city>New York</city>`,
			indent:      2,
			expectError: true,
		},
		{
			name:        "Same format - JSON to JSON",
			from:        "json",
			to:          "json",
			input:       `{"name":"John","age":30,"city":"New York"}`,
			indent:      2,
			expectError: false,
		},
		{
			name:        "YML alias to YAML",
			from:        "yml",
			to:          "yaml",
			input:       "name: John\nage: 30\ncity: New York",
			indent:      2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			converter := NewConverter(in, out)

			err := converter.Convert(tt.from, tt.to, tt.indent)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, out.String())
			}
		})
	}
}

func TestNewConverter(t *testing.T) {
	in := strings.NewReader("test")
	out := new(bytes.Buffer)

	converter := NewConverter(in, out)

	assert.NotNil(t, converter)
	assert.Equal(t, in, converter.reader)
	assert.Equal(t, out, converter.writer)
}

func TestConvertSameFormat(t *testing.T) {
	input := `{"name":"John","age":30,"city":"New York"}`
	in := strings.NewReader(input)
	out := new(bytes.Buffer)

	converter := NewConverter(in, out)
	err := converter.Convert("json", "json", 2)

	assert.NoError(t, err)
	assert.Equal(t, input, out.String())
}
