package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMainPackage - интеграционный тест для проверки, что package main компилируется и выполняется
func TestMainPackage(t *testing.T) {
	// Проверка наличия Go компилятора
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("Компилятор Go не найден, пропускаем тест")
	}

	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "devhelper-test")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Путь к временному исполняемому файлу
	binaryPath := tempDir + "/devhelper-test"

	// Компилируем
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Не удалось скомпилировать main.go: %v\nOutput: %s", err, output)
	}

	// Проверяем, что бинарный файл создан
	_, err = os.Stat(binaryPath)
	assert.NoError(t, err, "Бинарный файл должен существовать")

	// Запускаем с флагом --help
	cmd = exec.Command(binaryPath, "--help")
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Исполняемый файл должен успешно запускаться с флагом --help")

	// Проверяем, что вывод содержит ожидаемую информацию
	outputStr := string(output)
	assert.Contains(t, outputStr, "DevHelper", "Вывод должен содержать название приложения")
	assert.Contains(t, outputStr, "Usage:", "Вывод должен содержать информацию об использовании")
	assert.Contains(t, outputStr, "Available Commands:", "Вывод должен содержать список команд")
}

// TestVersionVariables проверяет, что переменные версии определены
func TestVersionVariables(t *testing.T) {
	assert.NotEmpty(t, Version, "Версия должна быть определена")

	// В тестовой среде BuildTime и GitCommit могут иметь значение "unknown"
	assert.NotNil(t, BuildTime, "BuildTime должен быть определен")
	assert.NotNil(t, GitCommit, "GitCommit должен быть определен")
}
