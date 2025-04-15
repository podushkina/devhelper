package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

// Config представляет конфигурацию приложения
type Config struct {
	// Общие настройки
	General struct {
		DefaultFormat string `json:"default_format" yaml:"default_format"`
		ColorEnabled  bool   `json:"color_enabled" yaml:"color_enabled"`
		DefaultIndent int    `json:"default_indent" yaml:"default_indent"`
	} `json:"general" yaml:"general"`

	// Настройки HTTP-клиента
	HTTP struct {
		Timeout           time.Duration `json:"timeout" yaml:"timeout"`
		FollowRedirects   bool          `json:"follow_redirects" yaml:"follow_redirects"`
		MaxRedirects      int           `json:"max_redirects" yaml:"max_redirects"`
		InsecureSSL       bool          `json:"insecure_ssl" yaml:"insecure_ssl"`
		DefaultUserAgent  string        `json:"default_user_agent" yaml:"default_user_agent"`
		DefaultHeaders    []string      `json:"default_headers" yaml:"default_headers"`
		SaveResponsesPath string        `json:"save_responses_path" yaml:"save_responses_path"`
	} `json:"http" yaml:"http"`

	// Настройки форматирования
	Formatter struct {
		JSONStyle  string `json:"json_style" yaml:"json_style"`
		YAMLStyle  string `json:"yaml_style" yaml:"yaml_style"`
		XMLStyle   string `json:"xml_style" yaml:"xml_style"`
		SortKeys   bool   `json:"sort_keys" yaml:"sort_keys"`
		WrapWidth  int    `json:"wrap_width" yaml:"wrap_width"`
		EscapeHTML bool   `json:"escape_html" yaml:"escape_html"`
	} `json:"formatter" yaml:"formatter"`

	// Настройки генератора
	Generator struct {
		DefaultCharset    string `json:"default_charset" yaml:"default_charset"`
		DefaultDateFormat string `json:"default_date_format" yaml:"default_date_format"`
		DefaultOutputType string `json:"default_output_type" yaml:"default_output_type"`
	} `json:"generator" yaml:"generator"`

	// Настройки монитора
	Monitor struct {
		DefaultInterval int    `json:"default_interval" yaml:"default_interval"`
		DefaultDisplay  string `json:"default_display" yaml:"default_display"`
		LogToFile       bool   `json:"log_to_file" yaml:"log_to_file"`
		LogFilePath     string `json:"log_file_path" yaml:"log_file_path"`
	} `json:"monitor" yaml:"monitor"`
}

// defaultConfig возвращает конфигурацию по умолчанию
func defaultConfig() *Config {
	var config Config

	// Общие настройки
	config.General.DefaultFormat = "json"
	config.General.ColorEnabled = true
	config.General.DefaultIndent = 2

	// HTTP-клиент
	config.HTTP.Timeout = 30 * time.Second
	config.HTTP.FollowRedirects = true
	config.HTTP.MaxRedirects = 10
	config.HTTP.InsecureSSL = false
	config.HTTP.DefaultUserAgent = "DevHelper/1.0"
	config.HTTP.DefaultHeaders = []string{}
	config.HTTP.SaveResponsesPath = ""

	// Форматирование
	config.Formatter.JSONStyle = "monokai"
	config.Formatter.YAMLStyle = "monokai"
	config.Formatter.XMLStyle = "monokai"
	config.Formatter.SortKeys = false
	config.Formatter.WrapWidth = 80
	config.Formatter.EscapeHTML = false

	// Генератор
	config.Generator.DefaultCharset = "alphanumeric"
	config.Generator.DefaultDateFormat = "2006-01-02"
	config.Generator.DefaultOutputType = "string"

	// Монитор
	config.Monitor.DefaultInterval = 1
	config.Monitor.DefaultDisplay = "dashboard"
	config.Monitor.LogToFile = false
	config.Monitor.LogFilePath = ""

	return &config
}

// ConfigManager управляет конфигурацией приложения
type ConfigManager struct {
	Config     *Config
	ConfigPath string
}

// NewConfigManager создает новый менеджер конфигурации
func NewConfigManager() (*ConfigManager, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию конфигурации: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	manager := &ConfigManager{
		Config:     defaultConfig(),
		ConfigPath: configPath,
	}

	// Загружаем существующую конфигурацию, если она есть
	if _, err := os.Stat(configPath); err == nil {
		if err := manager.Load(); err != nil {
			return nil, err
		}
	} else {
		// Если конфигурации нет, создаем ее
		if err := manager.Save(); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

// Load загружает конфигурацию из файла
func (m *ConfigManager) Load() error {
	data, err := os.ReadFile(m.ConfigPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	// Определяем формат файла на основе расширения
	ext := filepath.Ext(m.ConfigPath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, m.Config); err != nil {
			return fmt.Errorf("ошибка разбора JSON конфигурации: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, m.Config); err != nil {
			return fmt.Errorf("ошибка разбора YAML конфигурации: %w", err)
		}
	default:
		return fmt.Errorf("неподдерживаемое расширение файла конфигурации: %s", ext)
	}

	return nil
}

// Save сохраняет конфигурацию в файл
func (m *ConfigManager) Save() error {
	var data []byte
	var err error

	// Определяем формат файла на основе расширения
	ext := filepath.Ext(m.ConfigPath)
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(m.Config, "", "  ")
		if err != nil {
			return fmt.Errorf("ошибка сериализации JSON конфигурации: %w", err)
		}
	case ".yaml", ".yml":
		data, err = yaml.Marshal(m.Config)
		if err != nil {
			return fmt.Errorf("ошибка сериализации YAML конфигурации: %w", err)
		}
	default:
		return fmt.Errorf("неподдерживаемое расширение файла конфигурации: %s", ext)
	}

	if err := os.WriteFile(m.ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("не удалось записать файл конфигурации: %w", err)
	}

	return nil
}

// ResetToDefault сбрасывает конфигурацию к значениям по умолчанию
func (m *ConfigManager) ResetToDefault() error {
	m.Config = defaultConfig()
	return m.Save()
}

// getConfigDir возвращает путь к директории конфигурации
func getConfigDir() (string, error) {
	var configDir string

	// Определяем директорию в зависимости от ОС
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("не удалось определить APPDATA")
		}
		configDir = filepath.Join(appData, "DevHelper")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("не удалось определить домашнюю директорию: %w", err)
		}
		configDir = filepath.Join(home, "Library", "Application Support", "DevHelper")
	default: // Linux и другие UNIX-подобные ОС
		// Проверяем наличие XDG_CONFIG_HOME
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			configDir = filepath.Join(xdgConfigHome, "devhelper")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("не удалось определить домашнюю директорию: %w", err)
			}
			configDir = filepath.Join(home, ".config", "devhelper")
		}
	}

	return configDir, nil
}
