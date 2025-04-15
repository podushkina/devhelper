package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_GenerateUUID(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		format   string
		upper    bool
		validate func(t *testing.T, output string)
	}{
		{
			name:   "Single UUID - String format",
			count:  1,
			format: "string",
			upper:  false,
			validate: func(t *testing.T, output string) {
				// UUID v4 формат: 8-4-4-4-12 (всего 36 символов с дефисами)
				assert.Len(t, strings.TrimSpace(output), 36)
				assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, strings.TrimSpace(output))
			},
		},
		{
			name:   "Multiple UUIDs - String format",
			count:  3,
			format: "string",
			upper:  false,
			validate: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.Len(t, lines, 3)
				for _, line := range lines {
					assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, line)
				}
			},
		},
		{
			name:   "UUID with Uppercase - String format",
			count:  1,
			format: "string",
			upper:  true,
			validate: func(t *testing.T, output string) {
				assert.Len(t, strings.TrimSpace(output), 36)
				assert.Regexp(t, `^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`, strings.TrimSpace(output))
			},
		},
		{
			name:   "UUID in JSON format",
			count:  2,
			format: "json",
			upper:  false,
			validate: func(t *testing.T, output string) {
				var uuids []string
				err := json.Unmarshal([]byte(output), &uuids)
				assert.NoError(t, err)
				assert.Len(t, uuids, 2)
				for _, uuid := range uuids {
					assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, uuid)
				}
			},
		},
		{
			name:   "UUID in CSV format",
			count:  2,
			format: "csv",
			upper:  false,
			validate: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.Len(t, lines, 2)
				for _, line := range lines {
					// Проверка формата CSV (с кавычками)
					assert.True(t, strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\""))

					// Извлечение UUID из кавычек
					uuid := line[1 : len(line)-1]
					assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, uuid)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			generator := NewGenerator(out)

			err := generator.GenerateUUID(tt.count, tt.format, tt.upper)
			assert.NoError(t, err)
			tt.validate(t, out.String())
		})
	}
}

func TestGenerator_GenerateString(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		count    int
		charset  string
		format   string
		validate func(t *testing.T, output string, length int, charset string)
	}{
		{
			name:    "Alphanumeric String",
			length:  10,
			count:   1,
			charset: "alphanumeric",
			format:  "string",
			validate: func(t *testing.T, output string, length int, charset string) {
				assert.Len(t, strings.TrimSpace(output), length)
				assert.Regexp(t, `^[A-Za-z0-9]+$`, strings.TrimSpace(output))
			},
		},
		{
			name:    "Alpha String",
			length:  8,
			count:   1,
			charset: "alpha",
			format:  "string",
			validate: func(t *testing.T, output string, length int, charset string) {
				assert.Len(t, strings.TrimSpace(output), length)
				assert.Regexp(t, `^[A-Za-z]+$`, strings.TrimSpace(output))
			},
		},
		{
			name:    "Numeric String",
			length:  5,
			count:   1,
			charset: "numeric",
			format:  "string",
			validate: func(t *testing.T, output string, length int, charset string) {
				assert.Len(t, strings.TrimSpace(output), length)
				assert.Regexp(t, `^[0-9]+$`, strings.TrimSpace(output))
			},
		},
		{
			name:    "Hex String",
			length:  8,
			count:   1,
			charset: "hex",
			format:  "string",
			validate: func(t *testing.T, output string, length int, charset string) {
				assert.Len(t, strings.TrimSpace(output), length)
				assert.Regexp(t, `^[0-9a-f]+$`, strings.TrimSpace(output))
			},
		},
		{
			name:    "Multiple Strings - JSON format",
			length:  6,
			count:   3,
			charset: "alphanumeric",
			format:  "json",
			validate: func(t *testing.T, output string, length int, charset string) {
				var strings []string
				err := json.Unmarshal([]byte(output), &strings)
				assert.NoError(t, err)
				assert.Len(t, strings, 3)
				for _, s := range strings {
					assert.Len(t, s, length)
					assert.Regexp(t, `^[A-Za-z0-9]+$`, s)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			generator := NewGenerator(out)

			err := generator.GenerateString(tt.length, tt.count, tt.charset, tt.format)
			assert.NoError(t, err)
			tt.validate(t, out.String(), tt.length, tt.charset)
		})
	}

	// Проверка на неправильный charset
	t.Run("Invalid Charset", func(t *testing.T) {
		out := new(bytes.Buffer)
		generator := NewGenerator(out)

		err := generator.GenerateString(10, 1, "invalid", "string")
		assert.Error(t, err)
	})
}

func TestGenerator_GenerateNumber(t *testing.T) {
	tests := []struct {
		name     string
		min      int64
		max      int64
		count    int
		float    bool
		format   string
		validate func(t *testing.T, output string, min, max int64, isFloat bool)
	}{
		{
			name:   "Generate Integer",
			min:    1,
			max:    100,
			count:  1,
			float:  false,
			format: "string",
			validate: func(t *testing.T, output string, min, max int64, isFloat bool) {
				num := strings.TrimSpace(output)
				val := 0
				_, err := fmt.Sscanf(num, "%d", &val)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, int64(val), min)
				assert.LessOrEqual(t, int64(val), max)
			},
		},
		{
			name:   "Generate Float",
			min:    1,
			max:    100,
			count:  1,
			float:  true,
			format: "string",
			validate: func(t *testing.T, output string, min, max int64, isFloat bool) {
				num := strings.TrimSpace(output)
				val := 0.0
				_, err := fmt.Sscanf(num, "%f", &val)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, val, float64(min))
				assert.LessOrEqual(t, val, float64(max))
				// Проверка, что это действительно float (имеет десятичную часть)
				assert.Contains(t, num, ".")
			},
		},
		{
			name:   "Multiple Integers - JSON format",
			min:    1,
			max:    1000,
			count:  5,
			float:  false,
			format: "json",
			validate: func(t *testing.T, output string, min, max int64, isFloat bool) {
				var nums []string
				err := json.Unmarshal([]byte(output), &nums)
				assert.NoError(t, err)
				assert.Len(t, nums, 5)

				for _, num := range nums {
					val := 0
					_, err := fmt.Sscanf(num, "%d", &val)
					assert.NoError(t, err)
					assert.GreaterOrEqual(t, int64(val), min)
					assert.LessOrEqual(t, int64(val), max)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			generator := NewGenerator(out)

			err := generator.GenerateNumber(tt.min, tt.max, tt.count, tt.float, tt.format)
			assert.NoError(t, err)
			tt.validate(t, out.String(), tt.min, tt.max, tt.float)
		})
	}
}

func TestGenerator_GenerateDate(t *testing.T) {
	tests := []struct {
		name        string
		startStr    string
		endStr      string
		count       int
		dateFormat  string
		outFormat   string
		expectError bool
		validate    func(t *testing.T, output string, start, end time.Time, dateFormat string, count int)
	}{
		{
			name:        "Single Date - Default Format",
			startStr:    "2000-01-01",
			endStr:      "2023-01-01",
			count:       1,
			dateFormat:  "2006-01-02",
			outFormat:   "string",
			expectError: false,
			validate: func(t *testing.T, output string, start, end time.Time, dateFormat string, count int) {
				dateStr := strings.TrimSpace(output)
				date, err := time.Parse(dateFormat, dateStr)
				assert.NoError(t, err)
				assert.True(t, date.Equal(start) || date.After(start))
				assert.True(t, date.Equal(end) || date.Before(end))
			},
		},
		{
			name:        "Multiple Dates - Custom Format",
			startStr:    "2000-01-01",
			endStr:      "2023-01-01",
			count:       3,
			dateFormat:  "02/01/2006",
			outFormat:   "string",
			expectError: false,
			validate: func(t *testing.T, output string, start, end time.Time, dateFormat string, count int) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.Len(t, lines, count)

				for _, line := range lines {
					date, err := time.Parse(dateFormat, line)
					assert.NoError(t, err)
					assert.True(t, date.Equal(start) || date.After(start))
					assert.True(t, date.Equal(end) || date.Before(end))
				}
			},
		},
		{
			name:        "Invalid Start Date",
			startStr:    "invalid",
			endStr:      "2023-01-01",
			count:       1,
			dateFormat:  "2006-01-02",
			outFormat:   "string",
			expectError: true,
			validate:    func(t *testing.T, output string, start, end time.Time, dateFormat string, count int) {},
		},
		{
			name:        "Invalid End Date",
			startStr:    "2000-01-01",
			endStr:      "invalid",
			count:       1,
			dateFormat:  "2006-01-02",
			outFormat:   "string",
			expectError: true,
			validate:    func(t *testing.T, output string, start, end time.Time, dateFormat string, count int) {},
		},
		{
			name:        "End Date Before Start Date",
			startStr:    "2023-01-01",
			endStr:      "2000-01-01",
			count:       1,
			dateFormat:  "2006-01-02",
			outFormat:   "string",
			expectError: true,
			validate:    func(t *testing.T, output string, start, end time.Time, dateFormat string, count int) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			generator := NewGenerator(out)

			start, startErr := time.Parse("2006-01-02", tt.startStr)
			end, endErr := time.Parse("2006-01-02", tt.endStr)

			err := generator.GenerateDate(tt.startStr, tt.endStr, tt.count, tt.dateFormat, tt.outFormat)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, startErr)
				require.NoError(t, endErr)
				assert.NoError(t, err)
				tt.validate(t, out.String(), start, end, tt.dateFormat, tt.count)
			}
		})
	}
}
