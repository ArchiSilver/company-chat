```markdown
# Настольный клиент — мультиплатформенная заготовка (Wails)

Цель: получит небольшой нативный мультиплатформенный настольный клиент для Company Chat, повторно используя существующий HTML/JS интерфейс и оборачивая его лёгким Go‑бэкендом. Быстрый и удобный вариант для этого репозитория — Wails (https://wails.io/): он интегрируется с Go, переиспользует веб‑UI и генерирует нативные приложения для macOS, Linux и Windows.

Почему Wails
- Переиспользует `scripts/ui.html` и JS без переписывания интерфейса.
- Плотная интеграция с Go: можно вызывать функции Go из JS и упаковывать статические ресурсы.
- Кроссплатформенные сборки через `wails build` (macOS/Linux/Windows).

Основные шаги (локально)
1. Установите зависимости:
   - Go (рекомендуется 1.18+)
   - Node.js + npm (инструменты фронтенда)
   - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest` (убедитесь, что `$GOPATH/bin` в PATH)

2. Инициализация Wails проекта (однократно):
   ```sh
   mkdir -p desktop/wails && cd desktop/wails
   wails init -n company-chat-desktop -t vanilla
   ```

3. Замените фронтенд на существующий UI
   - Скопируйте `scripts/ui.html` (и ассеты) в фронтенд Wails (`frontend/src`) или адаптируйте index.html в шаблоне.
   - Обновите API_BASE → `http://localhost:8080` либо сделайте его настраиваемым через binding.

4. Запуск в режиме разработки
   - В `desktop/wails`: `wails dev` — откроется нативное окно с вашим UI, которое может обращаться к бэкенду на `localhost:8080`.

5. Сборка для дистрибуции
   - `wails build` создаст нативные артефакты: `.app` на macOS, `.exe` на Windows, бинарник на Linux.
   - Для кросс‑сборки может потребоваться дополнительный toolchain; в CI проще запускать сборку на соответствующих раннерах.

Заметки и рекомендации
- Для компактного мультиплатформенного дистрибутива Wails — хороший компромисс; альтернативой является Tauri (на Rust) с меньшими бинарями, но с более сложной настройкой.
- Если нужен полностью нативный UI на Go — рассмотрите `Fyne`, но придётся переписать интерфейс.
- Для оффлайн‑режима можно добавить встроенную SQLite и запускать локальный сервер вместе с десктопным приложением — это более трудоёмко и выполняется после базовой интеграции.

Makefile помощники
- В Makefile репозитория есть цели `make desktop-dev` и `make desktop-build`, которые вызывают Wails CLI (при наличии).

CI / Кросс‑сборки
- Для мультиплатформенных релизов используйте CI (GitHub Actions) с раннерами для macOS/Windows/Linux и вызывайте `wails build` на каждой платформе. Могу добавить пример workflow.

Дальше я могу:
- Сгенерировать заготовку `desktop/wails` и скопировать UI (локально потребует выполнить `wails init`).
- Добавить CI workflow для сборки desktop‑артефактов.
- Предложить план на Tauri вместо Wails.

```
# Desktop client — multiplatform (Wails)

Goal: produce a small, native, cross-platform desktop client for Company Chat by reusing the existing HTML/JS UI and building a lightweight Go wrapper. The fastest, most developer-friendly option for this repository is Wails (https://wails.io/) — it integrates with Go, reuses web UI, and builds native apps for macOS, Linux and Windows.

Why Wails
- Reuse your `scripts/ui.html` and JS without rewriting UI code.
- Tight Go integration: call Go functions from JS if needed, bundle static assets.
- Cross-platform builds via `wails build` (macOS/Linux/Windows).

High-level steps (what you'll run locally)
1. Install prerequisites
   - Go (you already have Go 1.24)
   - Node.js + npm (for Wails frontend tooling)
   - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest` (make sure `$GOPATH/bin` is in PATH)

2. Create a Wails project scaffold (one-time)
   - From repository root: `mkdir -p desktop/wails && cd desktop/wails`
   - `wails init -n company-chat-desktop -t vanilla`  # choose the vanilla template

3. Replace frontend with existing UI
   - Copy `scripts/ui.html` (and any assets) into the Wails frontend `frontend/src` (for vanilla template adjust index.html accordingly), or adapt frontend files per Wails template.
   - Update frontend code to point API_BASE to `http://localhost:8080` or make it configurable via the Wails binding.

4. Development run
   - From `desktop/wails`: `wails dev` — opens a live desktop window that uses your UI and can call backend endpoints on localhost:8080.

5. Build for distribution
   - `wails build` produces native artifacts. On macOS it will create a `.app`, on Windows an `.exe`, and on Linux a binary.
   - Cross-building may require extra toolchain (e.g., building Windows on macOS may need `xgo` or a CI runner on Windows). For local builds, run `wails build` on each target OS or use CI.

Notes & recommendations
- For a small-footprint multi-platform distribution, Wails is a good middle ground; Tauri is another option (Rust-based) with smaller bundles but requires a Rust toolchain and more setup.
- If you want fully native Go UI (no HTML), consider `Fyne` — but that requires rewriting the UI in Go widgets.
- Embedded/offline mode: if you want the desktop app to run without a remote server (embedded DB/server), we can add an embedded SQLite option and start a bundled server process from the desktop app — this is more work and can be done after the initial multiplatform client is ready.

Makefile helpers
- The repository Makefile contains helper targets `make desktop-dev` and `make desktop-build` that call the Wails CLI (they check for CLI presence first).

CI / Cross-platform builds
- For multi-platform releases build artifacts, use CI (GitHub Actions) with runners for macOS, Windows and Linux and invoke `wails build` on each runner. I can provide a workflow template.

If you'd like, I can:
- Scaffold `desktop/wails` with the minimal Wails project and copy the `scripts/ui.html` into the new frontend (requires wails CLI to complete some steps locally), or
- Create a CI workflow that builds desktop artifacts for macOS/Linux/Windows in GitHub Actions.

Which do you want next?
- "Scaffold Wails here" — I'll create initial files and copy UI; note: final local `wails init` may be required on your machine to install platform files.
- "Create CI for builds" — I'll add a GitHub Actions workflow that runs `wails build` on macOS/Windows/Linux runners.
- "Start with Tauri" — I'll propose a Tauri plan instead.
