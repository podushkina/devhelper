package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_SendRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		statusCode     int
		responseBody   string
		requestHeaders map[string]string
		requestBody    []byte
		username       string
		password       string
		insecure       bool
		expectError    bool
	}{
		{
			name:           "GET request",
			method:         "GET",
			statusCode:     http.StatusOK,
			responseBody:   `{"status":"ok"}`,
			requestHeaders: map[string]string{"Accept": "application/json"},
			requestBody:    nil,
			username:       "",
			password:       "",
			insecure:       false,
			expectError:    false,
		},
		{
			name:           "POST request with body",
			method:         "POST",
			statusCode:     http.StatusCreated,
			responseBody:   `{"id":1,"status":"created"}`,
			requestHeaders: map[string]string{"Content-Type": "application/json"},
			requestBody:    []byte(`{"name":"test"}`),
			username:       "",
			password:       "",
			insecure:       false,
			expectError:    false,
		},
		{
			name:           "Basic authentication",
			method:         "GET",
			statusCode:     http.StatusOK,
			responseBody:   `{"status":"authenticated"}`,
			requestHeaders: map[string]string{},
			requestBody:    nil,
			username:       "user",
			password:       "pass",
			insecure:       false,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый сервер
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				assert.Equal(t, tt.method, r.Method)

				// Проверяем заголовки
				for key, value := range tt.requestHeaders {
					assert.Equal(t, value, r.Header.Get(key))
				}

				// Проверяем базовую аутентификацию
				if tt.username != "" {
					username, password, ok := r.BasicAuth()
					assert.True(t, ok)
					assert.Equal(t, tt.username, username)
					assert.Equal(t, tt.password, password)
				}

				// Устанавливаем статус код и отправляем ответ
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Создаем HTTP-клиент
			client := NewHTTPClient(5 * time.Second)

			// Отправляем запрос
			response, err := client.SendRequest(
				tt.method,
				server.URL,
				tt.requestHeaders,
				tt.requestBody,
				tt.username,
				tt.password,
				tt.insecure,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, response.Status, fmt.Sprintf("%d", tt.statusCode))
				assert.Equal(t, tt.responseBody, string(response.Body))
				assert.NotZero(t, response.TotalTime)
			}
		})
	}
}

func TestHTTPClient_SendRequest_Error(t *testing.T) {
	// Создаем клиент с очень маленьким таймаутом
	client := NewHTTPClient(1 * time.Millisecond)

	// Создаем тестовый сервер, который будет "спать" дольше таймаута
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Спим дольше таймаута
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Отправляем запрос
	_, err := client.SendRequest(
		"GET",
		server.URL,
		nil,
		nil,
		"",
		"",
		false,
	)

	// Ожидаем ошибку таймаута
	assert.Error(t, err)
}

func TestHTTPClient_InvalidURL(t *testing.T) {
	client := NewHTTPClient(5 * time.Second)

	// Отправляем запрос с некорректным URL
	_, err := client.SendRequest(
		"GET",
		"http://invalid-url-that-does-not-exist.example",
		nil,
		nil,
		"",
		"",
		false,
	)

	// Ожидаем ошибку
	assert.Error(t, err)
}

func TestNewHTTPClient(t *testing.T) {
	timeout := 10 * time.Second
	client := NewHTTPClient(timeout)

	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
	assert.Equal(t, timeout, client.client.Timeout)
}
