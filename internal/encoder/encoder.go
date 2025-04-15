package encoder

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Encoder представляет кодировщик/декодировщик данных
type Encoder struct {
	reader io.Reader
	writer io.Writer
}

// NewEncoder создает новый кодировщик/декодировщик данных
func NewEncoder(reader io.Reader, writer io.Writer) *Encoder {
	return &Encoder{
		reader: reader,
		writer: writer,
	}
}

// NewCommand создает новую команду кодирования/декодирования
func NewCommand() *cobra.Command {
	encodeCmd := &cobra.Command{
		Use:   "encode",
		Short: "Кодирование и декодирование данных",
		Long:  "Кодирование и декодирование данных в различных форматах (Base64, URL).",
	}

	// Подкоманда для Base64
	base64Cmd := &cobra.Command{
		Use:   "base64",
		Short: "Кодирование/декодирование Base64",
		Long:  "Кодирование и декодирование данных в формате Base64.",
	}

	// Команда кодирования Base64
	base64EncodeCmd := &cobra.Command{
		Use:   "encode [string]",
		Short: "Кодирование в Base64",
		Long:  "Кодирование строки или файла в формат Base64.",
		Run: func(cmd *cobra.Command, args []string) {
			var input io.Reader

			if len(args) > 0 {
				input = strings.NewReader(args[0])
			} else {
				input = os.Stdin
			}

			urlSafe, _ := cmd.Flags().GetBool("urlsafe")
			encoder := NewEncoder(input, os.Stdout)
			if err := encoder.Base64Encode(urlSafe); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка кодирования в Base64: %s\n", err)
				os.Exit(1)
			}
		},
	}
	base64EncodeCmd.Flags().Bool("urlsafe", false, "Использовать URL-безопасный вариант Base64")

	// Команда декодирования Base64
	base64DecodeCmd := &cobra.Command{
		Use:   "decode [string]",
		Short: "Декодирование из Base64",
		Long:  "Декодирование строки из формата Base64.",
		Run: func(cmd *cobra.Command, args []string) {
			var input io.Reader

			if len(args) > 0 {
				input = strings.NewReader(args[0])
			} else {
				input = os.Stdin
			}

			urlSafe, _ := cmd.Flags().GetBool("urlsafe")
			encoder := NewEncoder(input, os.Stdout)
			if err := encoder.Base64Decode(urlSafe); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка декодирования из Base64: %s\n", err)
				os.Exit(1)
			}
		},
	}
	base64DecodeCmd.Flags().Bool("urlsafe", false, "Использовать URL-безопасный вариант Base64")

	// Подкоманда для URL
	urlCmd := &cobra.Command{
		Use:   "url",
		Short: "Кодирование/декодирование URL",
		Long:  "Кодирование и декодирование строк в URL-формат.",
	}

	// Команда кодирования URL
	urlEncodeCmd := &cobra.Command{
		Use:   "encode [string]",
		Short: "Кодирование в URL-формат",
		Long:  "Кодирование строки в URL-формат.",
		Run: func(cmd *cobra.Command, args []string) {
			var input io.Reader

			if len(args) > 0 {
				input = strings.NewReader(args[0])
			} else {
				input = os.Stdin
			}

			encoder := NewEncoder(input, os.Stdout)
			if err := encoder.URLEncode(); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка кодирования URL: %s\n", err)
				os.Exit(1)
			}
		},
	}

	// Команда декодирования URL
	urlDecodeCmd := &cobra.Command{
		Use:   "decode [string]",
		Short: "Декодирование из URL-формата",
		Long:  "Декодирование строки из URL-формата.",
		Run: func(cmd *cobra.Command, args []string) {
			var input io.Reader

			if len(args) > 0 {
				input = strings.NewReader(args[0])
			} else {
				input = os.Stdin
			}

			encoder := NewEncoder(input, os.Stdout)
			if err := encoder.URLDecode(); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка декодирования URL: %s\n", err)
				os.Exit(1)
			}
		},
	}

	// Добавление подкоманд
	base64Cmd.AddCommand(base64EncodeCmd, base64DecodeCmd)
	urlCmd.AddCommand(urlEncodeCmd, urlDecodeCmd)
	encodeCmd.AddCommand(base64Cmd, urlCmd)

	return encodeCmd
}

// Base64Encode кодирует данные в Base64
func (e *Encoder) Base64Encode(urlSafe bool) error {
	data, err := io.ReadAll(e.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	var encoded string
	if urlSafe {
		encoded = base64.URLEncoding.EncodeToString(data)
	} else {
		encoded = base64.StdEncoding.EncodeToString(data)
	}

	if _, err := fmt.Fprintln(e.writer, encoded); err != nil {
		return fmt.Errorf("ошибка записи результата: %w", err)
	}

	return nil
}

// Base64Decode декодирует данные из Base64
func (e *Encoder) Base64Decode(urlSafe bool) error {
	data, err := io.ReadAll(e.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	// Удаляем возможные пробелы и переводы строк
	cleanData := strings.TrimSpace(string(data))

	var decoded []byte
	if urlSafe {
		decoded, err = base64.URLEncoding.DecodeString(cleanData)
	} else {
		decoded, err = base64.StdEncoding.DecodeString(cleanData)
	}

	if err != nil {
		return fmt.Errorf("ошибка декодирования Base64: %w", err)
	}

	if _, err := e.writer.Write(decoded); err != nil {
		return fmt.Errorf("ошибка записи результата: %w", err)
	}

	return nil
}

// URLEncode кодирует данные в URL-формат
func (e *Encoder) URLEncode() error {
	data, err := io.ReadAll(e.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	encoded := url.QueryEscape(string(data))

	if _, err := fmt.Fprintln(e.writer, encoded); err != nil {
		return fmt.Errorf("ошибка записи результата: %w", err)
	}

	return nil
}

// URLDecode декодирует данные из URL-формата
func (e *Encoder) URLDecode() error {
	data, err := io.ReadAll(e.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	// Удаляем возможные пробелы и переводы строк
	cleanData := strings.TrimSpace(string(data))

	decoded, err := url.QueryUnescape(cleanData)
	if err != nil {
		return fmt.Errorf("ошибка декодирования URL: %w", err)
	}

	if _, err := fmt.Fprintln(e.writer, decoded); err != nil {
		return fmt.Errorf("ошибка записи результата: %w", err)
	}

	return nil
}
