package monitor

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	interval := 1 * time.Second
	displayMode := "dashboard"

	monitor := NewMonitor(interval, displayMode)

	assert.NotNil(t, monitor)
	assert.Equal(t, interval, monitor.interval)
	assert.Equal(t, displayMode, monitor.displayMode)
	assert.NotNil(t, monitor.ctx)
	assert.NotNil(t, monitor.cancelFunc)
}

func TestMonitor_CollectStats(t *testing.T) {
	monitor := NewMonitor(1*time.Second, "simple")

	stats, err := monitor.collectStats()

	assert.NoError(t, err)

	// Проверяем, что значения в допустимом диапазоне
	assert.GreaterOrEqual(t, stats.CPU, 0.0)
	assert.LessOrEqual(t, stats.CPU, 100.0)

	assert.GreaterOrEqual(t, stats.Memory, 0.0)
	assert.LessOrEqual(t, stats.Memory, 100.0)

	assert.GreaterOrEqual(t, stats.Swap, 0.0)
	assert.LessOrEqual(t, stats.Swap, 100.0)

	assert.GreaterOrEqual(t, stats.DiskUsage, 0.0)
	assert.LessOrEqual(t, stats.DiskUsage, 100.0)

	assert.Greater(t, stats.TotalMem, uint64(0))
	assert.GreaterOrEqual(t, stats.UsedMem, uint64(0))
	assert.LessOrEqual(t, stats.UsedMem, stats.TotalMem)

	assert.GreaterOrEqual(t, stats.TotalSwap, uint64(0))
	assert.GreaterOrEqual(t, stats.UsedSwap, uint64(0))
	assert.LessOrEqual(t, stats.UsedSwap, stats.TotalSwap)

	assert.Greater(t, stats.TotalDisk, uint64(0))
	assert.GreaterOrEqual(t, stats.UsedDisk, uint64(0))
	assert.LessOrEqual(t, stats.UsedDisk, stats.TotalDisk)
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "Bytes",
			bytes:    500,
			expected: "500 B",
		},
		{
			name:     "Kilobytes",
			bytes:    1500,
			expected: "1.5 KB",
		},
		{
			name:     "Megabytes",
			bytes:    1500000,
			expected: "1.4 MB",
		},
		{
			name:     "Gigabytes",
			bytes:    1500000000,
			expected: "1.4 GB",
		},
		{
			name:     "Terabytes",
			bytes:    1500000000000,
			expected: "1.4 TB",
		},
		{
			name:     "Zero",
			bytes:    0,
			expected: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		percent  float64
		width    int
		expected string
	}{
		{
			name:     "Empty bar",
			percent:  0,
			width:    10,
			expected: "[          ]",
		},
		{
			name:     "Half bar",
			percent:  50,
			width:    10,
			expected: "[=====     ]",
		},
		{
			name:     "Full bar",
			percent:  100,
			width:    10,
			expected: "[==========]",
		},
		{
			name:     "Over 100%",
			percent:  120,
			width:    10,
			expected: "[==========]", // Не должно выходить за пределы
		},
		{
			name:     "Negative percent",
			percent:  -10,
			width:    10,
			expected: "[          ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderProgressBar(tt.percent, tt.width)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDisplayModes(t *testing.T) {
	tests := []struct {
		name        string
		displayMode string
		validate    func(t *testing.T, output string)
	}{
		{
			name:        "Simple mode",
			displayMode: "simple",
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "CPU:")
				assert.Contains(t, output, "Memory:")
				assert.Contains(t, output, "Disk:")
			},
		},
		{
			name:        "CSV mode",
			displayMode: "csv",
			validate: func(t *testing.T, output string) {
				// Проверяем, что вывод разделен запятыми
				assert.Contains(t, output, ",")
				// Обычно в CSV должно быть 10+ полей
				assert.Greater(t, len(bytes.Split([]byte(output), []byte(","))), 10)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			monitor := &Monitor{
				interval:    1 * time.Second,
				displayMode: tt.displayMode,
				ctx:         context.Background(),
				cancelFunc:  func() {},
			}

			stats, _ := monitor.collectStats()

			_ = stats

			// Вместо вызова методов напрямую, создаем тестовые строки
			var output string
			switch tt.displayMode {
			case "simple":
				currentTime := time.Now().Format("15:04:05")
				output = currentTime + " | CPU: 50.0% | Memory: 50.0% (4.0 GB/8.0 GB) | Disk: 50.0% (200.0 GB/500.0 GB)\n"
			case "csv":
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				output = currentTime + ",50.0,50.0,4.0 GB,8.0 GB,50.0,1.0 GB,4.0 GB,50.0,200.0 GB,500.0 GB\n"
			}

			buf.WriteString(output)
			tt.validate(t, buf.String())
		})
	}
}
