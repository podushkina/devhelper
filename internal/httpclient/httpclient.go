package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// HTTPClient представляет HTTP-клиент
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient создает новый HTTP-клиент
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewCommand создает новую команду HTTP-клиента
func NewCommand() *cobra.Command {
	var (
		method      string
		headers     []string
		data        string
		dataFile    string
		timeout     int
		noColor     bool
		outputFile  string
		verbose     bool
		insecure    bool
		contentType string
		username    string
		password    string
		json        bool
	)

	httpCmd := &cobra.Command{
		Use:   "http [url]",
		Short: "HTTP-клиент для тестирования API",
		Long:  "Простой HTTP-клиент для отправки запросов и тестирования API.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]

			// Если не указан метод, используем GET
			if method == "" {
				method = "GET"
			}

			// Устанавливаем HTTP-клиент с таймаутом
			client := NewHTTPClient(time.Duration(timeout) * time.Second)

			// Если указан флаг --json, устанавливаем соответствующий Content-Type
			if json {
				contentType = "application/json"
			}

			// Собираем заголовки
			headerMap := make(map[string]string)
			for _, header := range headers {
				parts := strings.SplitN(header, ":", 2)
				if len(parts) == 2 {
					headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}

			// Если указан content-type, добавляем его в заголовки
			if contentType != "" {
				headerMap["Content-Type"] = contentType
			}

			// Определяем тело запроса
			var requestBody []byte
			if dataFile != "" {
				var err error
				requestBody, err = os.ReadFile(dataFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка чтения файла данных: %s\n", err)
					os.Exit(1)
				}
			} else if data != "" {
				requestBody = []byte(data)
			}

			// Отображаем спиннер во время запроса
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Выполнение запроса..."
			s.Start()

			// Выполняем запрос
			response, err := client.SendRequest(method, url, headerMap, requestBody, username, password, insecure)
			s.Stop()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при выполнении запроса: %s\n", err)
				os.Exit(1)
			}

			// Если указан выходной файл, сохраняем ответ в файл
			if outputFile != "" {
				if err := os.WriteFile(outputFile, response.Body, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка при сохранении ответа в файл: %s\n", err)
					os.Exit(1)
				}
				fmt.Printf("Ответ сохранен в файл: %s\n", outputFile)
				return
			}

			// Выводим информацию о запросе в вербозном режиме
			if verbose {
				// Вывод информации о запросе
				fmt.Printf("> %s %s\n", method, url)
				for key, value := range headerMap {
					fmt.Printf("> %s: %s\n", key, value)
				}
				if len(requestBody) > 0 {
					fmt.Println(">")
					fmt.Println(string(requestBody))
				}
				fmt.Println()
			}

			// Выводим информацию о статусе
			statusColor := color.New(color.FgCyan).SprintFunc()
			fmt.Printf("%s %s\n", statusColor(response.Status), response.Proto)

			// Выводим заголовки ответа
			printHeaders(response.Headers)

			// Выводим тело ответа с подсветкой синтаксиса, если это возможно
			printResponseBody(response.Body, response.Headers["Content-Type"], !noColor)
		},
	}

	httpCmd.Flags().StringVarP(&method, "method", "X", "", "HTTP-метод (GET, POST, PUT, DELETE и т.д.)")
	httpCmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "HTTP-заголовки (формат: 'Ключ: Значение')")
	httpCmd.Flags().StringVarP(&data, "data", "d", "", "Данные для отправки в теле запроса")
	httpCmd.Flags().StringVarP(&dataFile, "data-file", "f", "", "Файл с данными для отправки в теле запроса")
	httpCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Таймаут запроса в секундах")
	httpCmd.Flags().BoolVar(&noColor, "no-color", false, "Отключить подсветку синтаксиса")
	httpCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Сохранить ответ в файл")
	httpCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Подробный вывод")
	httpCmd.Flags().BoolVarP(&insecure, "insecure", "k", false, "Игнорировать проверку сертификатов SSL")
	httpCmd.Flags().StringVar(&contentType, "content-type", "", "Тип содержимого (Content-Type)")
	httpCmd.Flags().StringVarP(&username, "user", "u", "", "Имя пользователя и пароль для базовой аутентификации (формат: 'username:password')")
	httpCmd.Flags().StringVarP(&password, "password", "p", "", "Пароль для базовой аутентификации (если не указан в --user)")
	httpCmd.Flags().BoolVarP(&json, "json", "j", false, "Использовать Content-Type: application/json")

	return httpCmd
}

// HTTPResponse представляет ответ на HTTP-запрос
type HTTPResponse struct {
	Status    string
	Proto     string
	Headers   map[string]string
	Body      []byte
	TotalTime time.Duration
}

// SendRequest отправляет HTTP-запрос и возвращает ответ
func (c *HTTPClient) SendRequest(method, url string, headers map[string]string, body []byte, username, password string, insecure bool) (HTTPResponse, error) {
	startTime := time.Now()

	// Создаем запрос
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Добавляем заголовки
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Добавляем базовую аутентификацию, если указаны учетные данные
	if username != "" {
		auth := username
		_ = auth
		if password != "" {
			auth = fmt.Sprintf("%s:%s", username, password)
		}
		req.SetBasicAuth(username, password)
	}

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{}, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Собираем заголовки ответа
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		responseHeaders[key] = strings.Join(values, ", ")
	}

	return HTTPResponse{
		Status:    fmt.Sprintf("%d %s", resp.StatusCode, resp.Status),
		Proto:     resp.Proto,
		Headers:   responseHeaders,
		Body:      responseBody,
		TotalTime: time.Since(startTime),
	}, nil
}

// printHeaders выводит заголовки HTTP-ответа
func printHeaders(headers map[string]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Заголовок", "Значение"})

	// Сортируем заголовки для удобства чтения
	var sortedHeaders []string
	for key := range headers {
		sortedHeaders = append(sortedHeaders, key)
	}

	// Добавляем заголовки в таблицу
	for _, key := range sortedHeaders {
		t.AppendRow(table.Row{key, headers[key]})
	}

	t.SetStyle(table.StyleLight)
	t.Render()
	fmt.Println()
}

// printResponseBody выводит тело ответа с подсветкой синтаксиса
func printResponseBody(body []byte, contentType string, withColor bool) {
	if len(body) == 0 {
		return
	}

	// Пытаемся определить формат для подсветки
	var lexer chroma.Lexer

	// Пытаемся определить формат по Content-Type
	if contentType != "" {
		contentParts := strings.Split(contentType, ";")
		mimeType := strings.TrimSpace(contentParts[0])

		switch mimeType {
		case "application/json":
			lexer = lexers.Get("json")
		case "application/xml", "text/xml":
			lexer = lexers.Get("xml")
		case "text/html":
			lexer = lexers.Get("html")
		case "text/css":
			lexer = lexers.Get("css")
		case "application/javascript", "text/javascript":
			lexer = lexers.Get("javascript")
		}
	}

	// Если не удалось определить по Content-Type, пробуем определить по содержимому
	if lexer == nil {
		lexer = lexers.Analyse(string(body))
	}

	// Если и это не помогло, используем plaintext
	if lexer == nil {
		lexer = lexers.Get("plaintext")
	}

	// Проверяем, является ли содержимое валидным JSON для форматирования
	if lexer.Config().Name == "JSON" {
		var jsonObj interface{}
		var prettyJSON bytes.Buffer
		if err := json.Unmarshal(body, &jsonObj); err == nil {
			encoder := json.NewEncoder(&prettyJSON)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonObj); err == nil {
				body = prettyJSON.Bytes()
			}
		}
	}

	if withColor {
		// Подсветка синтаксиса
		style := styles.Get("monokai")
		formatter := formatters.Get("terminal")

		iterator, err := lexer.Tokenise(nil, string(body))
		if err != nil {
			fmt.Println(string(body))
			return
		}

		formatter.Format(os.Stdout, style, iterator)
	} else {
		fmt.Println(string(body))
	}
}
