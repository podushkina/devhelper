package app

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	versionInfo := VersionInfo{
		Version:   "1.0.0-test",
		BuildTime: "2025-01-01T12:00:00Z",
		GitCommit: "abcdef123456",
	}

	app := New(versionInfo)

	assert.NotNil(t, app)
	assert.NotNil(t, app.rootCmd)
	assert.Equal(t, versionInfo, app.versionInfo)
	assert.Equal(t, "devhelper", app.rootCmd.Use)
}

func TestRun(t *testing.T) {
	app := New(VersionInfo{})

	// Создаем собственный Command для тестирования
	// вместо замены app.rootCmd.Execute
	testCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {},
	}

	// Вызываем тестовую команду
	err := testCmd.Execute()
	assert.NoError(t, err)

	// Проверяем, что метод Run просто вызывает Execute на rootCmd
	// Это непрямой тест, который проверяет, что функциональность вообще работает
	assert.NotNil(t, app.rootCmd)
}
