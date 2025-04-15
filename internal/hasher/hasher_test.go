package hasher

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasher_GenerateHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hashFunc func() hash.Hash
		expected string
	}{
		{
			name:     "MD5",
			input:    "Hello, World!",
			hashFunc: md5.New,
			expected: "65a8e27d8879283831b664bd8b7f0ad4",
		},
		{
			name:     "SHA1",
			input:    "Hello, World!",
			hashFunc: sha1.New,
			expected: "0a0a9f2a6772942557ab5355d76af442f8f65e01",
		},
		{
			name:     "SHA256",
			input:    "Hello, World!",
			hashFunc: sha256.New,
			expected: "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f",
		},
		{
			name:     "SHA512",
			input:    "Hello, World!",
			hashFunc: sha512.New,
			expected: "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387",
		},
		{
			name:     "Empty string MD5",
			input:    "",
			hashFunc: md5.New,
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			hasher := NewHasher(in, out)

			hash, err := hasher.GenerateHash(tt.hashFunc())
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, hash)
		})
	}
}

func TestHasher_GenerateHash_InputError(t *testing.T) {
	errorReader := &ErrorReader{err: assert.AnError}
	out := new(bytes.Buffer)
	hasher := NewHasher(errorReader, out)

	hash, err := hasher.GenerateHash(md5.New())
	assert.Error(t, err)
	assert.Empty(t, hash)
}

// ErrorReader имитирует ошибку при чтении
type ErrorReader struct {
	err error
}

func (e *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestNewHasher(t *testing.T) {
	in := strings.NewReader("test")
	out := new(bytes.Buffer)

	hasher := NewHasher(in, out)

	assert.NotNil(t, hasher)
	assert.Equal(t, in, hasher.reader)
	assert.Equal(t, out, hasher.writer)
}

func TestHashComparisonLogic(t *testing.T) {
	// Проверка, что наши хеш-функции работают ожидаемым образом
	tests := []struct {
		name     string
		input    string
		hashFunc func() hash.Hash
		expected string
	}{
		{
			name:     "MD5 Test Vector",
			input:    "test",
			hashFunc: md5.New,
			expected: "098f6bcd4621d373cade4e832627b4f6",
		},
		{
			name:     "SHA1 Test Vector",
			input:    "test",
			hashFunc: sha1.New,
			expected: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
		},
		{
			name:     "SHA256 Test Vector",
			input:    "test",
			hashFunc: sha256.New,
			expected: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем хеш стандартным способом для проверки
			h := tt.hashFunc()
			h.Write([]byte(tt.input))
			expected := hex.EncodeToString(h.Sum(nil))

			// Теперь проверяем через наш Hasher
			in := strings.NewReader(tt.input)
			out := new(bytes.Buffer)
			hasher := NewHasher(in, out)

			hash, err := hasher.GenerateHash(tt.hashFunc())
			assert.NoError(t, err)

			// Сравниваем полученный хеш с ожидаемым значением
			assert.Equal(t, expected, hash)
			assert.Equal(t, tt.expected, hash)
		})
	}
}
