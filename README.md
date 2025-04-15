# DevHelper

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-green)

DevHelper - это универсальная многофункциональная консольная утилита для разработчиков, объединяющая набор необходимых инструментов в одном компактном приложении. Утилита значительно упрощает повседневные задачи, связанные с обработкой данных, форматированием, тестированием API и мониторингом ресурсов, повышая продуктивность без необходимости переключаться между различными инструментами.

## 🪄 Ключевые возможности

- **Форматирование данных** - высокопроизводительное форматирование и подсветка синтаксиса JSON, YAML и XML
- **Конвертация форматов** - беспрепятственное преобразование между JSON, YAML и XML с сохранением структуры
- **Генерация тестовых данных** - быстрое создание UUID, строк, чисел и дат для тестирования приложений
- **Кодирование/декодирование** - поддержка Base64 (стандартное и URL-safe) и URL кодирования
- **Вычисление хешей** - генерация и проверка хешей MD5, SHA1, SHA256 и SHA512
- **HTTP-клиент** - удобное тестирование API со всеми типами HTTP-запросов и заголовков
- **Мониторинг ресурсов** - наблюдение в реальном времени за использованием CPU, памяти и дисковой системы

### 🗺️ Планируемые функции

- Интеграция с системами автодополнения bash/zsh/fish
- Поддержка форматирования и валидации дополнительных форматов (TOML, CSV, INI)
- Генерация фиктивных данных по шаблонам (имена, адреса, телефоны)
- WebSocket клиент для тестирования в реальном времени
- Проверка сетевой доступности и DNS-поиск
- Интеграция с популярными фреймворками CI/CD
- Плагины расширения функциональности

## 📦 Содержание

- [Установка](#-установка)
- [Архитектура](#-архитектура)
- [Конфигурация](#-конфигурация)
- [Использование](#-использование)
    - [Форматирование данных](#форматирование-данных)
    - [Конвертация форматов](#конвертация-форматов)
    - [Генерация тестовых данных](#генерация-тестовых-данных)
    - [Кодирование/декодирование](#кодированиедекодирование)
    - [Вычисление хешей](#вычисление-хешей)
    - [HTTP-клиент](#http-клиент)
    - [Мониторинг ресурсов](#мониторинг-ресурсов)
- [Примеры](#-примеры)
- [Разработка](#-разработка)
- [Лицензия](#-лицензия)

## 🛠️ Установка

### Требования

- Go 1.24 или выше
- Доступ к интернету для скачивания зависимостей
- GNU Make (опционально)

### Установка из репозитория

```bash
# Клонирование репозитория
git clone https://github.com/podushkina/devhelper.git
cd devhelper

# Сборка проекта
make build

# Проверка работы
./build/devhelper version
```

### Установка в систему

```bash
# Использование make
make install

# Или ручная установка
go install ./cmd/devhelper
```

### Установка с использованием Go

```bash
go install github.com/podushkina/devhelper/cmd/devhelper@latest
```

### Скачивание бинарных файлов

Готовые бинарные файлы для различных платформ доступны на [странице релизов](https://github.com/podushkina/devhelper/releases).

> **Примечание для пользователей MacOS**: При запуске утилиты из dmg или zip файла может потребоваться дополнительное разрешение для выполнения. Используйте команду:
> ```bash
> chmod +x ./devhelper
> ```

## 🔍 Архитектура

DevHelper разработан с использованием модульной архитектуры, обеспечивающей высокую расширяемость и поддерживаемость кода.

### Компоненты системы

```
devhelper/
├── cmd/                  # Исполняемые приложения
│   └── devhelper/        # Основная точка входа
│       └── main.go
├── internal/             # Внутренний код приложения
│   ├── app/              # Основная логика приложения
│   │   └── app.go
│   ├── formatter/        # Форматирование JSON/YAML/XML
│   │   └── formatter.go
│   ├── converter/        # Конвертация между форматами
│   │   └── converter.go
│   ├── generator/        # Генерация тестовых данных
│   │   └── generator.go
│   ├── encoder/          # Кодирование/декодирование
│   │   └── encoder.go
│   ├── hasher/           # Генерация хэшей
│   │   └── hasher.go
│   ├── httpclient/       # HTTP-клиент
│   │   └── httpclient.go
│   └── monitor/          # Мониторинг ресурсов
│       └── monitor.go
├── pkg/                  # Публичный код библиотеки
│   ├── utils/            # Общие утилиты
│   │   └── utils.go
│   └── config/           # Конфигурация
│       └── config.go
├── tests/                # Тесты
│   └── integration/      # Интеграционные тесты
├── Makefile              # Инструкции для сборки
└── go.mod                # Описание модуля
```

### Принцип работы

1. **Загрузка конфигурации**: DevHelper загружает настройки из файла конфигурации или использует значения по умолчанию.
2. **Обработка команд**: Используется библиотека Cobra для обработки команд и флагов командной строки.
3. **Выполнение модулей**: Каждый модуль (formatter, converter и т.д.) отвечает за выполнение соответствующей функциональности.
4. **Вывод результатов**: Результаты выводятся в консоль или сохраняются в файл, в зависимости от настроек.

### Расширение функциональности

DevHelper спроектирован для легкого расширения:

1. **Добавление новых команд**: Создайте новый файл в директории `internal/` и реализуйте интерфейс Command.
2. **Настройка существующих команд**: Используйте систему конфигурации для тонкой настройки.
3. **Создание плагинов**: Следуйте шаблону плагина для расширения функциональности.

## ⚙️ Конфигурация

### Параметры командной строки

DevHelper предоставляет гибкие параметры конфигурации через командную строку. Общие флаги для всех команд:

| Параметр | Описание | Значение по умолчанию |
|----------|----------|------------------------|
| `--help` | Показать справку по команде | |
| `--no-color` | Отключить цветной вывод | `false` |
| `--output` | Путь к файлу для сохранения вывода | stdout |
| `--verbose` | Подробный вывод | `false` |

### Конфигурационный файл

DevHelper использует YAML-файл для расширенной конфигурации. По умолчанию ищет `.devhelper.yaml` в домашней директории пользователя.

```yaml
# Общие настройки
general:
  defaultFormat: "json"
  colorEnabled: true
  defaultIndent: 2

# Настройки HTTP-клиента
http:
  timeout: 30s
  followRedirects: true
  maxRedirects: 10
  insecureSSL: false
  defaultUserAgent: "DevHelper/1.0"
  defaultHeaders:
    - "Accept: application/json"
  saveResponsesPath: ""

# Настройки форматирования
formatter:
  jsonStyle: "monokai"
  yamlStyle: "monokai"
  xmlStyle: "monokai"
  sortKeys: false
  wrapWidth: 80
  escapeHTML: false

# Настройки генератора
generator:
  defaultCharset: "alphanumeric"
  defaultDateFormat: "2006-01-02"
  defaultOutputType: "string"

# Настройки монитора
monitor:
  defaultInterval: 1
  defaultDisplay: "dashboard"
  logToFile: false
  logFilePath: ""
```

### Переменные окружения

DevHelper также поддерживает конфигурацию через переменные окружения с префиксом `DEVHELPER_`:

| Переменная | Описание | Пример |
|------------|----------|--------|
| `DEVHELPER_CONFIG` | Путь к конфигурационному файлу | `~/.config/devhelper.yaml` |
| `DEVHELPER_COLOR` | Включить/выключить цветной вывод | `true` |
| `DEVHELPER_HTTP_TIMEOUT` | Таймаут HTTP-запросов в секундах | `30` |
| `DEVHELPER_HTTP_USER_AGENT` | User-Agent для HTTP-запросов | `MyApp/1.0` |

## 🚀 Использование

### Форматирование данных

```bash
# Форматирование JSON из файла
devhelper format json input.json

# Форматирование JSON с передачей через stdin
cat input.json | devhelper format json

# Форматирование YAML без подсветки синтаксиса
devhelper format yaml input.yaml --no-color

# Форматирование XML с пользовательским отступом
devhelper format xml input.xml --indent 4
```

Поддерживаемые опции:
- `--indent N` - установить размер отступа (по умолчанию 2)
- `--no-color` - отключить подсветку синтаксиса
- `--output FILE` - сохранить результат в файл

### Конвертация форматов

```bash
# Конвертация из JSON в YAML
devhelper convert json yaml input.json -o output.yaml

# Конвертация из YAML в JSON
devhelper convert yaml json input.yaml

# Конвертация из XML в JSON с настройкой отступа
devhelper convert xml json input.xml --indent 4
```

Поддерживаемые опции:
- `--indent N` - установить размер отступа (по умолчанию 2)
- `--output, -o FILE` - сохранить результат в файл

### Генерация тестовых данных

```bash
# Генерация одного UUID
devhelper generate uuid

# Генерация 5 UUID в формате JSON
devhelper generate uuid 5 --format json

# Генерация UUID в верхнем регистре
devhelper generate uuid --upper

# Генерация случайных строк
devhelper generate string 16 10 --charset alphanumeric

# Генерация случайных чисел
devhelper generate number 1 100 5 --format json

# Генерация случайных дат
devhelper generate date 2020-01-01 2023-12-31 3 --date-format "02.01.2006"
```

Поддерживаемые подкоманды:
- `uuid [count]` - генерация UUID
- `string [length] [count]` - генерация строк
- `number [min] [max] [count]` - генерация чисел
- `date [start] [end] [count]` - генерация дат

Опции:
- `--format, -f` - формат вывода (string, json, csv)
- `--charset, -c` - набор символов (alphanumeric, alpha, numeric, ascii, hex)
- `--upper, -u` - верхний регистр (для UUID)
- `--date-format, -d` - формат даты
- `--float` - генерировать числа с плавающей точкой

### Кодирование/декодирование

```bash
# Кодирование в Base64
devhelper encode base64 encode "Hello, World!"

# Декодирование из Base64
devhelper encode base64 decode "SGVsbG8sIFdvcmxkIQ=="

# Использование URL-безопасного Base64
devhelper encode base64 encode "Hello, World!+/" --urlsafe

# Кодирование URL
devhelper encode url encode "https://example.com/?query=test"

# Декодирование URL
devhelper encode url decode "https%3A%2F%2Fexample.com%2F%3Fquery%3Dtest"
```

Поддерживаемые подкоманды:
- `base64 encode|decode` - кодирование/декодирование Base64
- `url encode|decode` - кодирование/декодирование URL

Опции для Base64:
- `--urlsafe` - использовать URL-безопасный вариант Base64

### Вычисление хешей

```bash
# Вычисление MD5 хеша строки
devhelper hash md5 "Hello, World!"

# Вычисление SHA256 хеша файла
devhelper hash sha256 --file document.pdf

# Вычисление SHA1 с выводом в верхнем регистре
devhelper hash sha1 "test" --upper

# Проверка соответствия хеша
devhelper hash sha1 "test" --verify "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

# Тихая проверка хеша (только статус выхода)
devhelper hash sha256 --file image.iso --verify "a1b2c3..." --quiet
```

Поддерживаемые алгоритмы:
- `md5` - MD5 хеш
- `sha1` - SHA1 хеш
- `sha256` - SHA256 хеш
- `sha512` - SHA512 хеш

Опции:
- `--file, -f` - вычислить хеш файла вместо строки
- `--upper, -u` - вывести хеш в верхнем регистре
- `--verify, -v` - проверить совпадение хеша
- `--quiet, -q` - тихий режим (только статус выхода)

### HTTP-клиент

```bash
# GET запрос
devhelper http https://api.example.com/users

# POST запрос с JSON данными
devhelper http -X POST -H "Content-Type: application/json" -d '{"name":"John"}' https://api.example.com/users

# POST с данными из файла
devhelper http -X POST -f data.json https://api.example.com/users

# Подробный вывод с заголовками
devhelper http -v https://api.example.com/status

# Сохранение ответа в файл
devhelper http -o response.json https://api.example.com/data

# Использование аутентификации
devhelper http -u username:password https://api.example.com/secure

# Отключение проверки SSL
devhelper http -k https://self-signed.example.com
```

Опции:
- `--method, -X` - HTTP метод (GET, POST, PUT, DELETE и т.д.)
- `--header, -H` - HTTP заголовки
- `--data, -d` - данные для отправки в теле запроса
- `--data-file, -f` - файл с данными для отправки
- `--timeout, -t` - таймаут запроса в секундах
- `--output, -o` - сохранить ответ в файл
- `--verbose, -v` - подробный вывод
- `--insecure, -k` - игнорировать проверку сертификатов SSL
- `--user, -u` - имя пользователя и пароль для базовой аутентификации
- `--json, -j` - использовать Content-Type: application/json
- `--content-type` - тип содержимого (Content-Type)

### Мониторинг ресурсов

```bash
# Запуск мониторинга в режиме дашборда
devhelper monitor

# Изменение интервала обновления
devhelper monitor --interval 5

# Вывод в простом формате
devhelper monitor --display simple

# Вывод в формате CSV
devhelper monitor --display csv > metrics.csv
```

Опции:
- `--interval, -i` - интервал обновления в секундах
- `--display, -d` - режим отображения (dashboard, simple, csv)

## 📝 Примеры

### Пример обработки JSON данных

```bash
# Получение данных с API, форматирование и сохранение
devhelper http https://api.example.com/data | devhelper format json > formatted.json

# Преобразование формата и вычисление хеша
devhelper convert json yaml data.json | tee data.yaml | devhelper hash sha256
```

### Пример генерации тестовых данных для API

```bash
# Создание тестового JSON с UUID и отправка на API
echo '{"id":"'$(devhelper generate uuid)'","name":"Test"}' > data.json
devhelper http -X POST -f data.json https://api.example.com/create
```

### Пример вывода формата JSON

```json
{
  "name": "DevHelper",
  "version": "1.0.0",
  "features": [
    "Форматирование",
    "Конвертация",
    "Генерация данных",
    "Кодирование/декодирование",
    "Хеширование",
    "HTTP клиент",
    "Мониторинг"
  ],
  "configuration": {
    "defaultFormat": "json",
    "colorEnabled": true
  }
}
```

### Пример вывода монитора ресурсов

```
15:20:30 | DevHelper System Monitor

┌─────────┬───────────────────────────┬─────────┬──────────────────┐
│ Ресурс  │ Использование             │ Процент │ Детали           │
├─────────┼───────────────────────────┼─────────┼──────────────────┤
│ CPU     │ [==========          ]    │ 42.5%   │ 8 ядер           │
│ Memory  │ [=================   ]    │ 78.2%   │ 12.5 GB / 16 GB  │
│ Swap    │ [==                  ]    │ 11.3%   │ 1.1 GB / 8 GB    │
│ Disk    │ [===========         ]    │ 53.8%   │ 430 GB / 800 GB  │
└─────────┴───────────────────────────┴─────────┴──────────────────┘

Нажмите Ctrl+C для выхода
```

## 🛠️ Разработка

### Требования для разработки

- Go 1.24 или выше
- Редактор с поддержкой Go (VS Code, GoLand, Vim с плагинами)
- Git для контроля версий

### Сборка из исходного кода

```bash
# Клонирование репозитория
git clone https://github.com/podushkina/devhelper.git
cd devhelper

# Инициализация модуля
go mod tidy

# Запуск тестов
make test

# Сборка
make build
```

### Запуск тестов

```bash
# Запуск всех тестов
make test

# Запуск с покрытием кода
make cover

# Запуск линтера
make lint
```

### Кросс-компиляция

```bash
# Сборка для всех платформ
make build-all

# Сборка для конкретной платформы
make build-linux
make build-windows
make build-macos
```

### Вклад в проект

Вклады в проект приветствуются! Пожалуйста, ознакомьтесь с нашими рекомендациями по внесению вклада:

1. Форкните репозиторий и создайте новую ветку из `main`
2. Убедитесь, что код проходит все тесты (`make test`)
3. Добавьте новые тесты для новой функциональности
4. Отправьте Pull Request с подробным описанием изменений

## 📜 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для подробностей.

## 💌 Поддержка

Если у вас возникли вопросы или вы нашли ошибку, пожалуйста, создайте [issue](https://github.com/podushkina/devhelper/issues) в репозитории.

---

⭐️ Если вам нравится проект, не забудьте поставить ему звезду на GitHub! ⭐️
