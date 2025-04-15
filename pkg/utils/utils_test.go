package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "Microseconds",
			duration: 500 * time.Microsecond,
			expected: "500 µs",
		},
		{
			name:     "Milliseconds",
			duration: 500 * time.Millisecond,
			expected: "500 ms",
		},
		{
			name:     "Seconds",
			duration: 1500 * time.Millisecond,
			expected: "1.50 s",
		},
		{
			name:     "Minutes and seconds",
			duration: 90 * time.Second,
			expected: "1 min 30 s",
		},
		{
			name:     "Hours and minutes",
			duration: 90 * time.Minute,
			expected: "1 h 30 min",
		},
		{
			name:     "Zero",
			duration: 0,
			expected: "0 µs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "Bytes",
			bytes:    500,
			expected: "500 B",
		},
		{
			name:     "Kilobytes",
			bytes:    1500,
			expected: "1.5 KB",
		},
		{
			name:     "Megabytes",
			bytes:    1500000,
			expected: "1.4 MB",
		},
		{
			name:     "Gigabytes",
			bytes:    1500000000,
			expected: "1.4 GB",
		},
		{
			name:     "Terabytes",
			bytes:    1500000000000,
			expected: "1.4 TB",
		},
		{
			name:     "Zero",
			bytes:    0,
			expected: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileOperations(t *testing.T) {
	// Создаем временную директорию и файл для тестирования
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	assert.NoError(t, err)

	emptyDirPath := filepath.Join(tempDir, "emptydir")
	nonExistentPath := filepath.Join(tempDir, "nonexistent")

	// Тестируем FileExists
	t.Run("FileExists", func(t *testing.T) {
		assert.True(t, FileExists(testFilePath))
		assert.False(t, FileExists(nonExistentPath))
	})

	// Тестируем IsDirectory
	t.Run("IsDirectory", func(t *testing.T) {
		assert.True(t, IsDirectory(tempDir))
		assert.False(t, IsDirectory(testFilePath))
		assert.False(t, IsDirectory(nonExistentPath))
	})

	// Тестируем EnsureDir
	t.Run("EnsureDir", func(t *testing.T) {
		assert.NoError(t, EnsureDir(emptyDirPath))
		assert.True(t, IsDirectory(emptyDirPath))

		// Повторный вызов EnsureDir не должен давать ошибку
		assert.NoError(t, EnsureDir(emptyDirPath))
	})

	// Тестируем CopyFile
	t.Run("CopyFile", func(t *testing.T) {
		destFilePath := filepath.Join(tempDir, "destfile.txt")

		err := CopyFile(testFilePath, destFilePath)
		assert.NoError(t, err)

		// Проверяем, что содержимое скопировано корректно
		content, err := os.ReadFile(destFilePath)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		// Тестируем ошибку при копировании несуществующего файла
		err = CopyFile(nonExistentPath, destFilePath)
		assert.Error(t, err)
	})
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "String shorter than max",
			input:    "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "String equal to max",
			input:    "1234567890",
			maxLen:   10,
			expected: "1234567890",
		},
		{
			name:     "String longer than max",
			input:    "This is a very long string that should be truncated",
			maxLen:   20,
			expected: "This is a very lo...",
		},
		{
			name:     "Very short max",
			input:    "Hello",
			maxLen:   3,
			expected: "...",
		},
		{
			name:     "Empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetTerminalSize тестирует получение размеров терминала
// Этот тест может быть хрупким, так как зависит от окружения
func TestTerminalFunctions(t *testing.T) {
	// В автоматизированной среде тестирования stdin может не быть терминалом
	// Просто проверяем, что функция не паникует
	t.Run("IsTerminal", func(t *testing.T) {
		isTerminal := IsTerminal(int(os.Stdin.Fd()))
		// Не проверяем конкретное значение, просто убеждаемся что функция работает
		t.Logf("Is stdin a terminal: %v", isTerminal)
	})

	t.Run("GetTerminalSize", func(t *testing.T) {
		width, height, err := GetTerminalSize()
		// В CI среде может не быть реального терминала, поэтому ошибка допустима
		if err == nil {
			assert.GreaterOrEqual(t, width, 0)
			assert.GreaterOrEqual(t, height, 0)
		}
	})
}
