package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Converter представляет конвертер форматов данных
type Converter struct {
	reader io.Reader
	writer io.Writer
}

// NewConverter создает новый конвертер форматов данных
func NewConverter(reader io.Reader, writer io.Writer) *Converter {
	return &Converter{
		reader: reader,
		writer: writer,
	}
}

// NewCommand создает новую команду конвертации
func NewCommand() *cobra.Command {
	var (
		outputFile string
		indent     int
	)

	convertCmd := &cobra.Command{
		Use:   "convert [from] [to] [file]",
		Short: "Конвертация между форматами данных",
		Long: `Конвертация между различными форматами данных (JSON, YAML, XML).
Если файл не указан, ввод будет считан из stdin.`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			from := strings.ToLower(args[0])
			to := strings.ToLower(args[1])

			// Проверка поддерживаемых форматов
			supportedFormats := map[string]bool{
				"json": true,
				"yaml": true,
				"yml":  true,
				"xml":  true,
			}

			if !supportedFormats[from] {
				fmt.Fprintf(os.Stderr, "Неподдерживаемый исходный формат: %s\n", from)
				os.Exit(1)
			}

			if !supportedFormats[to] {
				fmt.Fprintf(os.Stderr, "Неподдерживаемый целевой формат: %s\n", to)
				os.Exit(1)
			}

			var input io.Reader = os.Stdin
			var output io.Writer = os.Stdout

			// Если указан входной файл
			if len(args) > 2 {
				file, err := os.Open(args[2])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка при открытии файла: %s\n", err)
					os.Exit(1)
				}
				defer file.Close()
				input = file
			}

			// Если указан выходной файл
			if outputFile != "" {
				file, err := os.Create(outputFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка при создании выходного файла: %s\n", err)
					os.Exit(1)
				}
				defer file.Close()
				output = file
			}

			converter := NewConverter(input, output)
			if err := converter.Convert(from, to, indent); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка конвертации: %s\n", err)
				os.Exit(1)
			}

			if outputFile != "" {
				absPath, _ := filepath.Abs(outputFile)
				fmt.Printf("Конвертация завершена, результат сохранен в: %s\n", absPath)
			}
		},
	}

	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Выходной файл (по умолчанию: stdout)")
	convertCmd.Flags().IntVar(&indent, "indent", 2, "Размер отступа для форматирования")

	return convertCmd
}

// Convert конвертирует данные из одного формата в другой
func (c *Converter) Convert(from, to string, indent int) error {
	// Нормализация форматов
	if from == "yml" {
		from = "yaml"
	}
	if to == "yml" {
		to = "yaml"
	}

	// Если форматы одинаковые, просто копируем данные
	if from == to {
		if _, err := io.Copy(c.writer, c.reader); err != nil {
			return fmt.Errorf("ошибка копирования данных: %w", err)
		}
		return nil
	}

	// Чтение входных данных
	data, err := io.ReadAll(c.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения входных данных: %w", err)
	}

	// Парсинг входных данных в промежуточную структуру
	var obj interface{}
	switch from {
	case "json":
		if err := json.Unmarshal(data, &obj); err != nil {
			return fmt.Errorf("ошибка разбора JSON: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(data, &obj); err != nil {
			return fmt.Errorf("ошибка разбора YAML: %w", err)
		}
	case "xml":
		if err := xml.Unmarshal(data, &obj); err != nil {
			return fmt.Errorf("ошибка разбора XML: %w", err)
		}
	}

	// Маршалинг в целевой формат
	var output []byte
	switch to {
	case "json":
		jsonIndent := strings.Repeat(" ", indent)
		output, err = json.MarshalIndent(obj, "", jsonIndent)
		if err != nil {
			return fmt.Errorf("ошибка сериализации в JSON: %w", err)
		}
	case "yaml":
		output, err = yaml.Marshal(obj)
		if err != nil {
			return fmt.Errorf("ошибка сериализации в YAML: %w", err)
		}
	case "xml":
		xmlIndent := strings.Repeat(" ", indent)
		output, err = xml.MarshalIndent(obj, "", xmlIndent)
		if err != nil {
			return fmt.Errorf("ошибка сериализации в XML: %w", err)
		}
	}

	// Запись результата
	_, err = c.writer.Write(output)
	if err != nil {
		return fmt.Errorf("ошибка записи результата: %w", err)
	}

	return nil
}
