package monitor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// SystemStats представляет статистику системных ресурсов
type SystemStats struct {
	CPU       float64 // Использование CPU (%)
	Memory    float64 // Использование памяти (%)
	UsedMem   uint64  // Использовано памяти (байты)
	TotalMem  uint64  // Всего памяти (байты)
	Swap      float64 // Использование swap (%)
	UsedSwap  uint64  // Использовано swap (байты)
	TotalSwap uint64  // Всего swap (байты)
	DiskUsage float64 // Использование диска (%)
	UsedDisk  uint64  // Использовано диска (байты)
	TotalDisk uint64  // Всего диска (байты)
}

// Monitor представляет монитор системных ресурсов
type Monitor struct {
	interval    time.Duration
	ctx         context.Context
	cancelFunc  context.CancelFunc
	displayMode string
}

// NewMonitor создает новый монитор системных ресурсов
func NewMonitor(interval time.Duration, displayMode string) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		interval:    interval,
		ctx:         ctx,
		cancelFunc:  cancel,
		displayMode: displayMode,
	}
}

// NewCommand создает новую команду мониторинга ресурсов
func NewCommand() *cobra.Command {
	var (
		interval    int
		displayMode string
	)

	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Мониторинг системных ресурсов",
		Long:  "Мониторинг использования CPU, памяти и диска.",
		Run: func(cmd *cobra.Command, args []string) {
			monitor := NewMonitor(time.Duration(interval)*time.Second, displayMode)
			if err := monitor.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при запуске мониторинга: %s\n", err)
				os.Exit(1)
			}
		},
	}

	monitorCmd.Flags().IntVarP(&interval, "interval", "i", 1, "Интервал обновления в секундах")
	monitorCmd.Flags().StringVarP(&displayMode, "display", "d", "dashboard", "Режим отображения (dashboard, simple, csv)")

	return monitorCmd
}

// Start запускает мониторинг системных ресурсов
func (m *Monitor) Start() error {
	// Обрабатываем сигналы для корректного завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем горутину для обработки сигналов
	go func() {
		<-sigCh
		m.cancelFunc()
	}()

	// Запускаем сбор статистики
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Если используем dashboard, очищаем экран и скрываем курсор
	if m.displayMode == "dashboard" {
		fmt.Print("\033[?25l")       // Скрываем курсор
		defer fmt.Print("\033[?25h") // Восстанавливаем курсор при выходе
	}

	// Если выводим CSV, печатаем заголовок
	if m.displayMode == "csv" {
		fmt.Println("Time,CPU (%),Memory (%),Memory Used,Memory Total,Swap (%),Swap Used,Swap Total,Disk (%),Disk Used,Disk Total")
	}

	for {
		select {
		case <-m.ctx.Done():
			return nil
		case <-ticker.C:
			stats, err := m.collectStats()
			if err != nil {
				return err
			}

			switch m.displayMode {
			case "dashboard":
				m.displayDashboard(stats)
			case "simple":
				m.displaySimple(stats)
			case "csv":
				m.displayCSV(stats)
			default:
				m.displayDashboard(stats)
			}
		}
	}
}

// collectStats собирает статистику системных ресурсов
func (m *Monitor) collectStats() (SystemStats, error) {
	var stats SystemStats
	var memStats runtime.MemStats

	// Заполняем статистику по памяти Go
	runtime.ReadMemStats(&memStats)

	// В реальном приложении здесь был бы код для сбора системных метрик
	// с использованием библиотек, например, github.com/shirou/gopsutil

	// Для примера используем фиктивные данные
	stats.CPU = 25.5
	stats.Memory = 60.2
	stats.UsedMem = 4 * 1024 * 1024 * 1024  // 4 GB
	stats.TotalMem = 8 * 1024 * 1024 * 1024 // 8 GB
	stats.Swap = 15.0
	stats.UsedSwap = 1 * 1024 * 1024 * 1024  // 1 GB
	stats.TotalSwap = 4 * 1024 * 1024 * 1024 // 4 GB
	stats.DiskUsage = 45.0
	stats.UsedDisk = 200 * 1024 * 1024 * 1024  // 200 GB
	stats.TotalDisk = 500 * 1024 * 1024 * 1024 // 500 GB

	return stats, nil
}

// displayDashboard отображает статистику в виде интерактивной панели
func (m *Monitor) displayDashboard(stats SystemStats) {
	// Очищаем экран
	fmt.Print("\033[H\033[2J")

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // Используем значение по умолчанию
	}

	// Выводим время
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n %s | DevHelper System Monitor\n\n", currentTime)

	// Создаем и настраиваем таблицу
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Ресурс", "Использование", "Процент", "Детали"})

	// Добавляем данные в таблицу
	cpuColor := getColorByPercent(stats.CPU)
	memColor := getColorByPercent(stats.Memory)
	swapColor := getColorByPercent(stats.Swap)
	diskColor := getColorByPercent(stats.DiskUsage)

	t.AppendRow(table.Row{
		"CPU",
		renderProgressBar(stats.CPU, width/3),
		cpuColor(fmt.Sprintf("%.1f%%", stats.CPU)),
		fmt.Sprintf("%d ядер", runtime.NumCPU()),
	})

	t.AppendRow(table.Row{
		"Memory",
		renderProgressBar(stats.Memory, width/3),
		memColor(fmt.Sprintf("%.1f%%", stats.Memory)),
		fmt.Sprintf("%s / %s", formatBytes(stats.UsedMem), formatBytes(stats.TotalMem)),
	})

	t.AppendRow(table.Row{
		"Swap",
		renderProgressBar(stats.Swap, width/3),
		swapColor(fmt.Sprintf("%.1f%%", stats.Swap)),
		fmt.Sprintf("%s / %s", formatBytes(stats.UsedSwap), formatBytes(stats.TotalSwap)),
	})

	t.AppendRow(table.Row{
		"Disk",
		renderProgressBar(stats.DiskUsage, width/3),
		diskColor(fmt.Sprintf("%.1f%%", stats.DiskUsage)),
		fmt.Sprintf("%s / %s", formatBytes(stats.UsedDisk), formatBytes(stats.TotalDisk)),
	})

	t.SetStyle(table.StyleLight)
	t.Render()

	fmt.Println("\nНажмите Ctrl+C для выхода")
}

// displaySimple отображает статистику в простом формате
func (m *Monitor) displaySimple(stats SystemStats) {
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("%s | CPU: %.1f%% | Memory: %.1f%% (%s/%s) | Disk: %.1f%% (%s/%s)\n",
		currentTime,
		stats.CPU,
		stats.Memory, formatBytes(stats.UsedMem), formatBytes(stats.TotalMem),
		stats.DiskUsage, formatBytes(stats.UsedDisk), formatBytes(stats.TotalDisk),
	)
}

// displayCSV отображает статистику в формате CSV
func (m *Monitor) displayCSV(stats SystemStats) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s,%.1f,%.1f,%s,%s,%.1f,%s,%s,%.1f,%s,%s\n",
		currentTime,
		stats.CPU,
		stats.Memory, formatBytes(stats.UsedMem), formatBytes(stats.TotalMem),
		stats.Swap, formatBytes(stats.UsedSwap), formatBytes(stats.TotalSwap),
		stats.DiskUsage, formatBytes(stats.UsedDisk), formatBytes(stats.TotalDisk),
	)
}

// renderProgressBar создает строку прогресс-бара
func renderProgressBar(percent float64, width int) string {
	if width < 10 {
		width = 10
	}

	completed := int(percent / 100 * float64(width))
	if completed > width {
		completed = width
	}

	bar := "["
	for i := 0; i < width; i++ {
		if i < completed {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"
	return bar
}

// formatBytes форматирует байты в читаемый формат
func formatBytes(bytes uint64) string {
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

// getColorByPercent возвращает функцию цвета в зависимости от процента
func getColorByPercent(percent float64) func(a ...interface{}) string {
	if percent >= 90 {
		return color.New(color.FgRed).SprintFunc()
	} else if percent >= 70 {
		return color.New(color.FgYellow).SprintFunc()
	}
	return color.New(color.FgGreen).SprintFunc()
}
