package generator

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Generator представляет генератор тестовых данных
type Generator struct {
	writer io.Writer
}

// NewGenerator создает новый генератор тестовых данных
func NewGenerator(writer io.Writer) *Generator {
	return &Generator{
		writer: writer,
	}
}

// NewCommand создает новую команду генерации тестовых данных
func NewCommand() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "generate",
		Short: "Генерация тестовых данных",
		Long:  "Генерация различных типов тестовых данных для разработки и тестирования.",
	}

	// Подкоманда для генерации UUID
	uuidCmd := &cobra.Command{
		Use:   "uuid [count]",
		Short: "Генерация UUID",
		Long:  "Генерация одного или нескольких UUID v4.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			count := 1
			var err error
			if len(args) > 0 {
				count, err = strconv.Atoi(args[0])
				if err != nil || count < 1 {
					fmt.Fprintf(os.Stderr, "Неверное количество: %s\n", args[0])
					os.Exit(1)
				}
			}

			format, _ := cmd.Flags().GetString("format")
			upperCase, _ := cmd.Flags().GetBool("upper")

			generator := NewGenerator(os.Stdout)
			if err := generator.GenerateUUID(count, format, upperCase); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при генерации UUID: %s\n", err)
				os.Exit(1)
			}
		},
	}
	uuidCmd.Flags().StringP("format", "f", "string", "Формат вывода (string, json, csv)")
	uuidCmd.Flags().BoolP("upper", "u", false, "Преобразовать UUID в верхний регистр")

	// Подкоманда для генерации строк
	stringCmd := &cobra.Command{
		Use:   "string [length] [count]",
		Short: "Генерация случайных строк",
		Long:  "Генерация одной или нескольких случайных строк заданной длины.",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			length := 10
			count := 1
			var err error

			if len(args) > 0 {
				length, err = strconv.Atoi(args[0])
				if err != nil || length < 1 {
					fmt.Fprintf(os.Stderr, "Неверная длина: %s\n", args[0])
					os.Exit(1)
				}
			}

			if len(args) > 1 {
				count, err = strconv.Atoi(args[1])
				if err != nil || count < 1 {
					fmt.Fprintf(os.Stderr, "Неверное количество: %s\n", args[1])
					os.Exit(1)
				}
			}

			charset, _ := cmd.Flags().GetString("charset")
			format, _ := cmd.Flags().GetString("format")

			generator := NewGenerator(os.Stdout)
			if err := generator.GenerateString(length, count, charset, format); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при генерации строк: %s\n", err)
				os.Exit(1)
			}
		},
	}
	stringCmd.Flags().StringP("charset", "c", "alphanumeric", "Набор символов (alphanumeric, alpha, numeric, ascii, hex)")
	stringCmd.Flags().StringP("format", "f", "string", "Формат вывода (string, json, csv)")

	// Подкоманда для генерации чисел
	numberCmd := &cobra.Command{
		Use:   "number [min] [max] [count]",
		Short: "Генерация случайных чисел",
		Long:  "Генерация одного или нескольких случайных чисел в заданном диапазоне.",
		Args:  cobra.MaximumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			min := int64(1)
			max := int64(100)
			count := 1
			var err error

			if len(args) > 0 {
				min, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Неверное минимальное значение: %s\n", args[0])
					os.Exit(1)
				}
			}

			if len(args) > 1 {
				max, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil || max < min {
					fmt.Fprintf(os.Stderr, "Неверное максимальное значение: %s\n", args[1])
					os.Exit(1)
				}
			}

			if len(args) > 2 {
				count, err = strconv.Atoi(args[2])
				if err != nil || count < 1 {
					fmt.Fprintf(os.Stderr, "Неверное количество: %s\n", args[2])
					os.Exit(1)
				}
			}

			format, _ := cmd.Flags().GetString("format")
			float, _ := cmd.Flags().GetBool("float")

			generator := NewGenerator(os.Stdout)
			if err := generator.GenerateNumber(min, max, count, float, format); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при генерации чисел: %s\n", err)
				os.Exit(1)
			}
		},
	}
	numberCmd.Flags().StringP("format", "f", "string", "Формат вывода (string, json, csv)")
	numberCmd.Flags().Bool("float", false, "Генерировать дробные числа вместо целых")

	// Подкоманда для генерации дат
	dateCmd := &cobra.Command{
		Use:   "date [start] [end] [count]",
		Short: "Генерация случайных дат",
		Long:  "Генерация одной или нескольких случайных дат в заданном диапазоне.",
		Args:  cobra.MaximumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			startStr := "2000-01-01"
			endStr := time.Now().Format("2006-01-02")
			count := 1
			var err error

			if len(args) > 0 {
				startStr = args[0]
			}

			if len(args) > 1 {
				endStr = args[1]
			}

			if len(args) > 2 {
				count, err = strconv.Atoi(args[2])
				if err != nil || count < 1 {
					fmt.Fprintf(os.Stderr, "Неверное количество: %s\n", args[2])
					os.Exit(1)
				}
			}

			format, _ := cmd.Flags().GetString("format")
			outputFormat, _ := cmd.Flags().GetString("date-format")

			generator := NewGenerator(os.Stdout)
			if err := generator.GenerateDate(startStr, endStr, count, outputFormat, format); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при генерации дат: %s\n", err)
				os.Exit(1)
			}
		},
	}
	dateCmd.Flags().StringP("format", "f", "string", "Формат вывода (string, json, csv)")
	dateCmd.Flags().StringP("date-format", "d", "2006-01-02", "Формат даты (Go time format)")

	// Добавляем подкоманды
	genCmd.AddCommand(uuidCmd)
	genCmd.AddCommand(stringCmd)
	genCmd.AddCommand(numberCmd)
	genCmd.AddCommand(dateCmd)

	return genCmd
}

// GenerateUUID генерирует UUID
func (g *Generator) GenerateUUID(count int, format string, upperCase bool) error {
	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		id := uuid.New().String()
		if upperCase {
			id = strings.ToUpper(id)
		}
		uuids[i] = id
	}

	return g.outputResult(uuids, format)
}

// GenerateString генерирует случайные строки
func (g *Generator) GenerateString(length, count int, charset, format string) error {
	var chars string
	switch charset {
	case "alphanumeric":
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	case "alpha":
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "numeric":
		chars = "0123456789"
	case "ascii":
		chars = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	case "hex":
		chars = "0123456789abcdef"
	default:
		return fmt.Errorf("неизвестный набор символов: %s", charset)
	}

	charCount := len(chars)
	strings := make([]string, count)

	for i := 0; i < count; i++ {
		result := make([]byte, length)
		for j := 0; j < length; j++ {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(charCount)))
			if err != nil {
				return fmt.Errorf("ошибка генерации случайного числа: %w", err)
			}
			result[j] = chars[n.Int64()]
		}
		strings[i] = string(result)
	}

	return g.outputResult(strings, format)
}

// GenerateNumber генерирует случайные числа
func (g *Generator) GenerateNumber(min, max int64, count int, float bool, format string) error {
	if float {
		floats := make([]float64, count)
		diff := float64(max - min)
		for i := 0; i < count; i++ {
			// Генерируем случайное число с помощью crypto/rand
			randBytes := make([]byte, 8)
			_, err := rand.Read(randBytes)
			if err != nil {
				return fmt.Errorf("ошибка генерации случайного числа: %w", err)
			}
			// Преобразуем байты в число от 0 до 1
			n := float64(int64(randBytes[0])|int64(randBytes[1])<<8|int64(randBytes[2])<<16|int64(randBytes[3])<<24) / float64(1<<32)
			if n < 0 {
				n = -n
			}
			if n > 1 {
				n = 1 / n
			}
			floats[i] = float64(min) + diff*n
		}

		// Преобразуем float64 в строки для вывода
		result := make([]string, count)
		for i, f := range floats {
			result[i] = fmt.Sprintf("%.6f", f)
		}
		return g.outputResult(result, format)
	}

	range64 := big.NewInt(max - min + 1)
	ints := make([]string, count)
	for i := 0; i < count; i++ {
		n, err := rand.Int(rand.Reader, range64)
		if err != nil {
			return fmt.Errorf("ошибка генерации случайного числа: %w", err)
		}
		ints[i] = fmt.Sprintf("%d", n.Int64()+min)
	}

	return g.outputResult(ints, format)
}

// GenerateDate генерирует случайные даты
func (g *Generator) GenerateDate(startStr, endStr string, count int, dateFormat, outputFormat string) error {
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return fmt.Errorf("неверный формат начальной даты: %w", err)
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return fmt.Errorf("неверный формат конечной даты: %w", err)
	}

	if end.Before(start) {
		return fmt.Errorf("конечная дата должна быть после начальной даты")
	}

	diff := end.Sub(start)
	diffDays := int64(diff.Hours() / 24)

	dates := make([]string, count)
	for i := 0; i < count; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(diffDays+1))
		if err != nil {
			return fmt.Errorf("ошибка генерации случайного числа: %w", err)
		}

		randomDate := start.Add(time.Duration(n.Int64()) * 24 * time.Hour)
		dates[i] = randomDate.Format(dateFormat)
	}

	return g.outputResult(dates, outputFormat)
}

// outputResult выводит результат в заданном формате
func (g *Generator) outputResult(data []string, format string) error {
	switch format {
	case "string":
		for _, item := range data {
			if _, err := fmt.Fprintln(g.writer, item); err != nil {
				return err
			}
		}
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("ошибка сериализации в JSON: %w", err)
		}
		if _, err := g.writer.Write(jsonData); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(g.writer); err != nil {
			return err
		}
	case "csv":
		for _, item := range data {
			if _, err := fmt.Fprintf(g.writer, "\"%s\"\n", item); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("неизвестный формат вывода: %s", format)
	}

	return nil
}
