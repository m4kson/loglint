# loglint

Линтер для проверки лог‑сообщений в Go‑проектах.  
Поддерживает `log/slog`, `go.uber.org/zap` и стандартный `log`, реализован как `go/analysis`‑анализатор и интегрируется в `golangci-lint` через **Module Plugin System**.

---

## Содержание

- [Возможности](#возможности)
- [Поддерживаемые логгеры](#поддерживаемые-логгеры)
- [Правила проверки](#правила-проверки)
  - [1. Первая буква — строчная](#1-первая-буква--строчная)
  - [2. Только английский (ASCII)](#2-только-английский-ascii)
  - [3. Без спецсимволов и эмодзи](#3-без-спецсимволов-и-эмодзи)
  - [4. Без потенциально чувствительных данных](#4-без-потенциально-чувствительных-данных)
- [Установка](#установка)
  - [Требования](#требования)
  - [Установка как обычного Go‑модуля](#установка-как-обычного-goмодуля)
- [Использование](#использование)
  - [1. Стендэлон‑анализатор (singlechecker)](#1-стендэлонанализатор-singlechecker)
  - [2. Интеграция с golangci-lint (Module Plugin System)](#2-интеграция-с-golangci-lint-module-plugin-system)
    - [2.1. Быстрый старт для стороннего проекта](#21-быстрый-старт-для-стороннего-проекта)
    - [2.2. Пример `.custom-gcl.yml` в стороннем проекте](#22-пример-custom-gclyml-в-стороннем-проекте)
    - [2.3. Пример `.golangci.yml` в стороннем проекте](#23-пример-golangciyml-в-стороннем-проекте)
- [Ограничения и особенности](#ограничения-и-особенности)
- [Разработка и вклад](#разработка-и-вклад)
  - [Локальный запуск линтера в этом репозитории](#локальный-запуск-линтера-в-этом-репозитории)
  - [Структура репозитория](#структура-репозитория)

---

## Возможности

- **Анализ лог‑вызовов** с использованием `go/analysis`.
- **Поддержка нескольких логгеров**: `log/slog`, `go.uber.org/zap`, стандартный `log`.
- **Набор правил для логов**:
  - стиль сообщений (первая буква),
  - язык (английский),
  - отсутствие спецсимволов и эмодзи,
  - защита от логирования потенциально чувствительных данных.
- **SuggestedFix** для некоторых правил (например, автоматическое исправление первой буквы на строчную).
- **Интеграция с `golangci-lint` как кастомный линтер** через Module Plugin System.

---

## Поддерживаемые логгеры

**log/slog** (`log/slog`):

- Методы: `Debug`, `Info`, `Warn`, `Error`, а также `DebugContext`, `InfoContext`, `WarnContext`, `ErrorContext`.
- Сообщение — первый аргумент (`msg` или форматная строка).

**go.uber.org/zap**:

- Методы базового логгера: `Debug`, `Info`, `Warn`, `Error`, `Fatal`, `Panic`, `DPanic`.
- Sugared‑логгер:
  - форматные методы: `Debugf`, `Infof`, `Warnf`, `Errorf`, `Fatalf`,
  - структурные методы: `Debugw`, `Infow`, `Warnw`, `Errorw`, `Fatalw`.
- Сообщение — первый аргумент.

**log** (stdlib):

- Методы: `Print`, `Println`, `Printf`, `Fatal`, `Fatalf`, `Fatalln`, `Panic`, `Panicf`, `Panicln`.
- Анализируется строковый литерал, если он первым аргументом и является строкой.

---

## Правила проверки

### 1. Первая буква — строчная

**Неправильно:**

```go
slog.Info("Server started")
slog.Error("Failed to connect")
```

**Правильно:**

```go
slog.Info("server started")
slog.Error("failed to connect")
```

- Линтер проверяет **первую Unicode‑руну** сообщения.
- Если это заглавная буква (`unicode.IsUpper`), будет диагностировано нарушение:

> `log message must start with a lowercase letter`

- Для этого правила реализован **SuggestedFix**:
  - golangci-lint может автоматически заменить первую букву на строчную.

### 2. Только английский

Сообщение должно содержать **только печатные ASCII‑символы** в диапазоне `0x20`–`0x7E`.

**Нарушения:**

- Кириллица: `"сервер запущен"`,
- Латинские буквы с диакритикой: `"café"`, `"naïve"`,
- Любые другие non‑ASCII символы.

Пример:

```go
slog.Info("сервер запущен") // нарушение
slog.Info("café au lait")   // нарушение
slog.Info("server started") // ок
```

Сообщение:

> `log message must contain only English (ASCII) characters, found '…'`

### 3. Без спецсимволов и эмодзи

Цель — не допускать в логах посторонних символов и эмодзи, которые мешают чтению и анализу.

- Для ASCII‑символов разрешены:
  - буквы,
  - цифры,
  - пробел,
  - символ `%` (часто используется в форматных строках).
- Все остальные ASCII‑символы (`#`, `!`, `?`, многоточия и т.п.) считаются нарушением.
- Для non‑ASCII‑символов используется заранее собранная `unicode.RangeTable` с диапазонами эмодзи и графических символов.

Примеры нарушений:

```go
slog.Info("connection failed!!!")             // спецсимволы
slog.Warn("warning: something went wrong...") // спецсимволы
slog.Error("this is #3")                      // спецсимвол '#'
slog.Info("😀 well done")                     // эмодзи
```

Сообщения:

- для спецсимволов:  
  `log message must not contain special characters, found '#'`
- для эмодзи:  
  `log message must not contain emoji, found '😀'`

### 4. Без потенциально чувствительных данных

Линтер ищет в тексте сообщения **подозрительные ключевые слова**:

- `"password"`, `"passwd"`, `"pwd"`,
- `"secret"`, `"token"`, `"bearer"`, `"jwt"`,
- `"apikey"`, `"api_key"`, `"api-key"`,
- `"auth"`, `"credential"`,
- `"private_key"`, `"privatekey"`, `"access_key"`, `"accesskey"`,
- `"ssn"`, `"credit_card"`, `"creditcard"`, `"cvv"`, `"pin"`,
- и др.

Поиск идёт по нижнему регистру всего сообщения (`strings.ToLower` + `strings.Contains`).

Примеры:

```go
slog.Info("user password reset")  // нарушение
slog.Error("token expired")       // нарушение
slog.Info("invalid jwt provided") // нарушение
slog.Info("bearer token missing") // нарушение
```

Сообщение:

> `log message may contain sensitive data: keyword "password" found in message`

---

## Установка

### Требования

- **Go**: 1.22+ (проект собирается под 1.25.x).
- **golangci-lint**: 2.x (для использования Module Plugin System).

### Установка как обычного Go‑модуля

В стороннем проекте:

```bash
go get github.com/m4kson/loglint@latest
```

Это позволит:

- использовать встроенный анализатор `analyzer.New()` в своих инструментах,
- подключить модуль как плагин для `golangci-lint` (см. ниже).

---

## Использование

### 1. singlechecker

В этом репозитории есть бинарный файл `cmd/loglint`, который запускает анализатор как `singlechecker`.

В корне проекта:

```bash
go run ./cmd/loglint ./...
```

или, после установки:

```bash
go install github.com/m4kson/loglint/cmd/loglint@latest

loglint ./...
```

Это запустит `go/analysis`‑анализатор по переданным пакетам и выведет диагностические сообщения по описанным выше правилам.

---

### 2. Интеграция с golangci-lint (Module Plugin System)

Рекомендуемый способ использования линтера в сторонних проектах - как **module‑plugin** для `golangci-lint`.

#### 2.1. Быстрый старт для стороннего проекта

В вашем **стороннем проекте**:

1. **Установите golangci-lint** (если его ещё нет):

   ```bash
   go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
   ```

2. **Добавьте зависимость на `loglint`**:

   ```bash
   go get github.com/m4kson/loglint@latest
   ```

3. **Создайте файл `.custom-gcl.yml`** (в корне вашего проекта):

   См. пример ниже.

4. **Создайте/обновите `.golangci.yml`**, чтобы:
   - включить кастомный линтер `loglint`,
   - описать его как `type: "module"`.

5. **Соберите кастомный golangci-lint**:

   ```bash
   golangci-lint custom -v
   ```

   По умолчанию будет создан бинарник `./custom-gcl` (или то имя, что вы укажете в `.custom-gcl.yml`).

6. **Запустите линтинг вашего проекта**:

   ```bash
   ./custom-gcl run ./... --config=.golangci.yml
   ```

#### 2.2. Пример `.custom-gcl.yml` в стороннем проекте

Минимальный пример:

```yaml
version: v2.1.5          # версия golangci-lint, с которой вы хотите собрать кастомный бинарник
name: custom-gcl         # имя итогового бинаря
destination: ./bin       # куда положить бинарь

plugins:
  - module: github.com/m4kson/loglint
    # Если вы используете loglint как обычную go-зависимость из Go proxy,
    # достаточно module + version. В большинстве случаев import не нужен,
    # так как плагин регистрируется в подпакете plugin.
    # Если вы хотите фиксированную версию:
    # version: vX.Y.Z
```

> В этом репозитории пример `.custom-gcl.yml` находится в корне и собирает кастомный бинарник для самого проекта.

#### 2.3. Пример `.golangci.yml` в стороннем проекте

Минимальная конфигурация, включающая только `loglint`:

```yaml
version: "2"

run:
  timeout: 5m
  relative-path-mode: gomod
  issues-exit-code: 1
  tests: true
  modules-download-mode: mod

issues:
  # Разрешаем несколько ошибок на одной строке от одного и того же линтера
  # (иначе golangci-lint может схлопывать репорты до одной проблемы на строку).
  uniq-by-line: false

linters:
  default: none
  enable:
    - loglint

  settings:
    custom:
      loglint:
        type: "module"
        description: "Checks Go log messages for style, language and sensitive data"
        settings: {}
```

После этого:

```bash
# сборка кастомного golangci-lint с плагином loglint
golangci-lint custom -v

# запуск линтера (предположим, бинарь называется ./bin/custom-gcl)
./bin/custom-gcl run ./... --config=.golangci.yml
```

---

## Ограничения и особенности

- **Анализируются только строковые литералы** в позициях `msg`/форматной строки:
  - выражения вида `someLogger.Info(msg)` с переменной `msg` **не анализируются** (это осознанный компромисс).
- **Проверка английского языка** реализована как ограничение до печатного ASCII:
  - все non‑ASCII символы (включая корректные английские буквы с диакритикой) считаются нарушением;
  - это упрощает имплементацию и делает правило более строгим.
- **Спецсимволы** определены консервативно:
  - лог‑сообщения ориентированы на человекочитаемый текст и анализ логов,
  - символы вроде `#`, `!`, `?`, `...` расцениваются как нежелательные.
- **Чувствительные данные** ищутся по простым подстрокам в нижнем регистре:
  - возможны как **ложноположительные**, так и **ложноотрицательные** срабатывания,
  - правило скорее предупреждающее (signal), чем полностью формально доказуемое.

---

## Разработка и вклад

### Локальный запуск линтера в этом репозитории

**Только анализатор (go/analysis + тесты):**

```bash
go test ./...
```

**Запуск встроенного бинарика `cmd/loglint`:**

```bash
go run ./cmd/loglint ./...
```

**Запуск через golangci-lint с модульным плагином (для самого проекта):**

```bash
# в корне этого репозитория
task lint
```

Команда:

- установит `golangci-lint` (если нужно),
- соберёт кастомный `./bin/custom-gcl` с модульным плагином `loglint`,
- запустит `./bin/custom-gcl run ./... --config=.golangci.yml`.

### Структура репозитория

- `cmd/loglint` — бинарь на основе `singlechecker`, запускающий анализатор.
- `pkg/analyzer` — основной анализатор:
  - `analyzer.go` — фабрика `New()` и основной `Run`,
  - `detector/detector.go` — поиск вызовов логгеров и извлечение сообщений,
  - `rules/*.go` — реализации правил (`lowercase`, `english-only`, `no-special-chars`, `no-sensitive-data`),
  - `testdata/` — тестовые кейсы для `analysistest`.
- `plugin/loglint_plugin.go` — реализация module‑плагина для `golangci-lint`:
  - регистрация через `github.com/golangci/plugin-module-register/register`,
  - `BuildAnalyzers()` возвращает `analyzer.New()`.
- `.custom-gcl.yml` — конфигурация сборки кастомного `golangci-lint` (для этого репозитория).
- `.golangci.yml` — минимальный конфиг, включающий только `loglint`.
- `Taskfile.yml` — удобные команды для:
  - установки `golangci-lint`,
  - сборки `custom-gcl`,
  - запуска `lint`.

---

