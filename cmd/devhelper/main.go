package main

import (
	"fmt"
	"os"

	"devhelper/internal/app"
)

// Version информация о версии приложения
var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Создаем экземпляр приложения с информацией о версии
	application := app.New(app.VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	})

	// Запускаем приложение и обрабатываем возможные ошибки
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %s\n", err)
		os.Exit(1)
	}
}
