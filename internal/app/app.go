package app

import (
	"fmt"
	"os"

	"devhelper/internal/converter"
	"devhelper/internal/encoder"
	"devhelper/internal/formatter"
	"devhelper/internal/generator"
	"devhelper/internal/hasher"
	"devhelper/internal/httpclient"
	"devhelper/internal/monitor"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// VersionInfo содержит информацию о версии приложения
type VersionInfo struct {
	Version   string
	BuildTime string
	GitCommit string
}

// App представляет основное приложение
type App struct {
	rootCmd     *cobra.Command
	versionInfo VersionInfo
}

// New создает новый экземпляр приложения
func New(versionInfo VersionInfo) *App {
	app := &App{
		versionInfo: versionInfo,
	}

	// Создаем корневую команду
	app.rootCmd = &cobra.Command{
		Use:   "devhelper",
		Short: "DevHelper - многофункциональный инструмент для разработчиков",
		Long: `DevHelper - многофункциональный инструмент для разработчиков,
который объединяет несколько полезных функций в одной CLI-утилите.

Возможности:
  * Форматирование JSON/YAML/XML с подсветкой синтаксиса
  * Конвертация между форматами данных
  * Генерация случайных тестовых данных
  * Кодирование/декодирование Base64, URL
  * Генерация хэшей (MD5, SHA1, SHA256)
  * Простой HTTP-клиент для тестирования API
  * Мониторинг использования системных ресурсов`,
		Run: func(cmd *cobra.Command, args []string) {
			// Если нет подкоманды, показываем справку
			cmd.Help()
		},
	}

	// Добавляем флаг версии
	app.rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Показать информацию о версии",
		Run: func(cmd *cobra.Command, args []string) {
			cyan := color.New(color.FgCyan).SprintFunc()
			fmt.Printf("%s: %s\n", cyan("Version"), app.versionInfo.Version)
			fmt.Printf("%s: %s\n", cyan("Build Time"), app.versionInfo.BuildTime)
			fmt.Printf("%s: %s\n", cyan("Git Commit"), app.versionInfo.GitCommit)
		},
	})

	// Регистрируем все команды
	app.registerCommands()

	return app
}

// Run запускает приложение
func (a *App) Run() error {
	return a.rootCmd.Execute()
}

// registerCommands регистрирует все команды приложения
func (a *App) registerCommands() {
	// Форматирование
	formatterCmd := formatter.NewCommand()
	a.rootCmd.AddCommand(formatterCmd)

	// Конвертация
	converterCmd := converter.NewCommand()
	a.rootCmd.AddCommand(converterCmd)

	// Генерация тестовых данных
	generatorCmd := generator.NewCommand()
	a.rootCmd.AddCommand(generatorCmd)

	// Кодирование/декодирование
	encoderCmd := encoder.NewCommand()
	a.rootCmd.AddCommand(encoderCmd)

	// Генерация хэшей
	hasherCmd := hasher.NewCommand()
	a.rootCmd.AddCommand(hasherCmd)

	// HTTP-клиент
	httpCmd := httpclient.NewCommand()
	a.rootCmd.AddCommand(httpCmd)

	// Мониторинг ресурсов
	monitorCmd := monitor.NewCommand()
	a.rootCmd.AddCommand(monitorCmd)

	// Добавляем команду для завершения shell
	a.rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Сгенерировать скрипт автодополнения для указанной оболочки",
		Long: `Сгенерировать скрипт автодополнения для devhelper для указанной оболочки.

Для bash:
  $ source <(devhelper completion bash)

Для zsh:
  $ source <(devhelper completion zsh)

Для fish:
  $ devhelper completion fish | source

Для powershell:
  PS> devhelper completion powershell | Out-String | Invoke-Expression
`,
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	})
}
