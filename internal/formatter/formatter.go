package formatter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

// Formatter представляет форматер данных
type Formatter struct {
	reader io.Reader
	writer io.Writer
}

// NewFormatter создает новый форматер данных
func NewFormatter(reader io.Reader, writer io.Writer) *Formatter {
	return &Formatter{
		reader: reader,
		writer: writer,
	}
}

// NewCommand создает новую команду форматирования
func NewCommand() *cobra.Command {
	var (
		noColor bool
		indent  int
	)

	formatCmd := &cobra.Command{
		Use:   "format [json|yaml|xml] [file]",
		Short: "Форматирование и подсветка JSON/YAML/XML",
		Long: `Форматирование и подсветка синтаксиса для файлов JSON, YAML и XML.
Если файл не указан, ввод будет считан из stdin.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			format := args[0]
			var input io.Reader = os.Stdin

			if len(args) > 1 {
				file, err := os.Open(args[1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка при открытии файла: %s\n", err)
					os.Exit(1)
				}
				defer file.Close()
				input = file
			}

			formatter := NewFormatter(input, os.Stdout)
			var err error

			switch strings.ToLower(format) {
			case "json":
				err = formatter.FormatJSON(indent, !noColor)
			case "yaml", "yml":
				err = formatter.FormatYAML(!noColor)
			case "xml":
				err = formatter.FormatXML(indent, !noColor)
			default:
				fmt.Fprintf(os.Stderr, "Неподдерживаемый формат: %s\n", format)
				os.Exit(1)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка форматирования: %s\n", err)
				os.Exit(1)
			}
		},
	}

	formatCmd.Flags().BoolVar(&noColor, "no-color", false, "Отключить подсветку синтаксиса")
	formatCmd.Flags().IntVar(&indent, "indent", 2, "Размер отступа для форматирования")

	return formatCmd
}

// FormatJSON форматирует JSON
func (f *Formatter) FormatJSON(indent int, color bool) error {
	data, err := io.ReadAll(f.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	// Проверка валидности JSON
	var jsonObj interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return fmt.Errorf("невалидный JSON: %w", err)
	}

	// Форматирование с библиотекой pretty
	opts := &pretty.Options{
		Width:    80,
		Prefix:   "",
		Indent:   strings.Repeat(" ", indent),
		SortKeys: false,
	}
	formatted := pretty.PrettyOptions(data, opts)

	if color {
		// Подсветка синтаксиса с Chroma
		lexer := lexers.Get("json")
		style := styles.Get("monokai")
		formatter := formatters.Get("terminal")

		iterator, err := lexer.Tokenise(nil, string(formatted))
		if err != nil {
			return fmt.Errorf("ошибка токенизации: %w", err)
		}

		return formatter.Format(f.writer, style, iterator)
	}

	_, err = f.writer.Write(formatted)
	return err
}

// FormatYAML форматирует YAML
func (f *Formatter) FormatYAML(color bool) error {
	data, err := io.ReadAll(f.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	// Проверка валидности YAML
	var yamlObj interface{}
	if err := yaml.Unmarshal(data, &yamlObj); err != nil {
		return fmt.Errorf("невалидный YAML: %w", err)
	}

	// Форматирование
	formatted, err := yaml.Marshal(yamlObj)
	if err != nil {
		return fmt.Errorf("ошибка форматирования YAML: %w", err)
	}

	if color {
		// Подсветка синтаксиса с Chroma
		lexer := lexers.Get("yaml")
		style := styles.Get("monokai")
		formatter := formatters.Get("terminal")

		iterator, err := lexer.Tokenise(nil, string(formatted))
		if err != nil {
			return fmt.Errorf("ошибка токенизации: %w", err)
		}

		return formatter.Format(f.writer, style, iterator)
	}

	_, err = f.writer.Write(formatted)
	return err
}

// FormatXML форматирует XML
func (f *Formatter) FormatXML(indent int, color bool) error {
	data, err := io.ReadAll(f.reader)
	if err != nil {
		return fmt.Errorf("ошибка чтения данных: %w", err)
	}

	// Проверка валидности XML
	var xmlObj interface{}
	if err := xml.Unmarshal(data, &xmlObj); err != nil {
		return fmt.Errorf("невалидный XML: %w", err)
	}

	// Форматирование XML
	var formattedBuf strings.Builder
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	encoder := xml.NewEncoder(&formattedBuf)
	encoder.Indent("", strings.Repeat(" ", indent))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка декодирования XML: %w", err)
		}

		if err := encoder.EncodeToken(token); err != nil {
			return fmt.Errorf("ошибка кодирования XML: %w", err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("ошибка при завершении кодирования XML: %w", err)
	}

	formatted := []byte(formattedBuf.String())

	if color {
		// Подсветка синтаксиса с Chroma
		lexer := lexers.Get("xml")
		style := styles.Get("monokai")
		formatter := formatters.Get("terminal")

		iterator, err := lexer.Tokenise(nil, string(formatted))
		if err != nil {
			return fmt.Errorf("ошибка токенизации: %w", err)
		}

		return formatter.Format(f.writer, style, iterator)
	}

	_, err = f.writer.Write(formatted)
	return err
}
