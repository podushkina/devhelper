package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// IsTerminal проверяет, является ли файловый дескриптор терминалом
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// GetTerminalSize возвращает размеры терминала
func GetTerminalSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// FormatDuration форматирует продолжительность в удобочитаемый формат
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.2f s", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%d min %d s", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%d h %d min", int(d.Hours()), int(d.Minutes())%60)
}

// FormatBytes форматирует байты в удобочитаемый формат
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Confirm запрашивает у пользователя подтверждение (да/нет)
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// PrintError выводит сообщение об ошибке в stderr
func PrintError(format string, a ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s\n", red("ОШИБКА:"), fmt.Sprintf(format, a...))
}

// PrintWarning выводит предупреждение
func PrintWarning(format string, a ...interface{}) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Fprintf(os.Stderr, "%s %s\n", yellow("ПРЕДУПРЕЖДЕНИЕ:"), fmt.Sprintf(format, a...))
}

// PrintSuccess выводит сообщение об успехе
func PrintSuccess(format string, a ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s %s\n", green("УСПЕХ:"), fmt.Sprintf(format, a...))
}

// PrintInfo выводит информационное сообщение
func PrintInfo(format string, a ...interface{}) {
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Fprintf(os.Stdout, "%s %s\n", cyan("ИНФО:"), fmt.Sprintf(format, a...))
}

// CopyFile копирует файл из src в dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

// FileExists проверяет, существует ли файл
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory проверяет, является ли путь директорией
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir создает директорию, если она не существует
func EnsureDir(dir string) error {
	if !FileExists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// TruncateString обрезает строку, если она длиннее максимальной длины
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
