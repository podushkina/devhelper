package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Hasher представляет генератор хэшей
type Hasher struct {
	reader io.Reader
	writer io.Writer
}

// NewHasher создает новый генератор хэшей
func NewHasher(reader io.Reader, writer io.Writer) *Hasher {
	return &Hasher{
		reader: reader,
		writer: writer,
	}
}

// NewCommand создает новую команду генерации хэшей
func NewCommand() *cobra.Command {
	hashCmd := &cobra.Command{
		Use:   "hash [algorithm] [string]",
		Short: "Генерация хэшей",
		Long: `Генерация хэшей для заданной строки или файла с использованием различных алгоритмов.
Поддерживаемые алгоритмы: md5, sha1, sha256, sha512.
Если строка не указана, данные будут считаны из stdin.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			algorithm := strings.ToLower(args[0])

			var input io.Reader
			if len(args) > 1 {
				input = strings.NewReader(args[1])
			} else {
				input = os.Stdin
			}

			file, _ := cmd.Flags().GetString("file")
			if file != "" {
				f, err := os.Open(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Ошибка при открытии файла: %s\n", err)
					os.Exit(1)
				}
				defer f.Close()
				input = f
			}

			upper, _ := cmd.Flags().GetBool("upper")
			verify, _ := cmd.Flags().GetString("verify")
			quiet, _ := cmd.Flags().GetBool("quiet")

			hasher := NewHasher(input, os.Stdout)

			var hashValue string
			var err error

			switch algorithm {
			case "md5":
				hashValue, err = hasher.GenerateHash(md5.New())
			case "sha1":
				hashValue, err = hasher.GenerateHash(sha1.New())
			case "sha256":
				hashValue, err = hasher.GenerateHash(sha256.New())
			case "sha512":
				hashValue, err = hasher.GenerateHash(sha512.New())
			default:
				fmt.Fprintf(os.Stderr, "Неподдерживаемый алгоритм: %s\n", algorithm)
				os.Exit(1)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при генерации хэша: %s\n", err)
				os.Exit(1)
			}

			if upper {
				hashValue = strings.ToUpper(hashValue)
			}

			// Если указана проверка хэша
			if verify != "" {
				if !quiet {
					fmt.Fprintf(os.Stderr, "Проверка хэша: ")
				}

				if strings.EqualFold(hashValue, verify) {
					// Хэши совпадают
					if !quiet {
						color.Green("OK\n")
					}
					os.Exit(0)
				} else {
					// Хэши не совпадают
					if !quiet {
						color.Red("ОШИБКА\n")
						fmt.Fprintf(os.Stderr, "Ожидается: %s\n", verify)
						fmt.Fprintf(os.Stderr, "Получено: %s\n", hashValue)
					}
					os.Exit(1)
				}
			}

			// Если не задана проверка, просто выводим хэш
			fmt.Fprintln(os.Stdout, hashValue)
		},
	}

	hashCmd.Flags().StringP("file", "f", "", "Использовать файл вместо строки или stdin")
	hashCmd.Flags().BoolP("upper", "u", false, "Вывести хэш в верхнем регистре")
	hashCmd.Flags().StringP("verify", "v", "", "Проверить совпадение хэша с указанным значением")
	hashCmd.Flags().BoolP("quiet", "q", false, "Тихий режим (без вывода сообщений, только статус выхода)")

	return hashCmd
}

// GenerateHash генерирует хэш с использованием указанного алгоритма
func (h *Hasher) GenerateHash(hashFunc hash.Hash) (string, error) {
	data, err := io.ReadAll(h.reader)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения данных: %w", err)
	}

	hashFunc.Write(data)
	hashValue := hex.EncodeToString(hashFunc.Sum(nil))

	return hashValue, nil
}
