```markdown
# EventyOn — заметки по интеграции с Wails

Эта папка содержит лёгкую заготовку для интеграции UI в проект на Wails. Wails CLI генерирует несколько платформенно‑зависимых файлов, поэтому в репозитории оставлен placeholder и инструкции для завершения настройки локально.

Шаги для завершения и сборки Wails‑приложения (локально):

1. Установите зависимости
   - Go (1.18+)
   - Node.js + npm
   - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2. Инициализируйте Wails проект в `desktop/wails`:

```sh
cd desktop/wails
wails init -n eventyon-desktop -t vanilla
```

Это создаст `frontend/` и `backend/`. Скопируйте содержимое `scripts/tg_ui.html` в фронтенд (или адаптируйте шаблон) чтобы переиспользовать UI.

3. Запуск в режиме разработки:

```sh
cd desktop/wails
wails dev
```

Откроется нативное окно с вашим UI. Убедитесь, что бэкенд доступен по `http://localhost:8080`.

4. Сборка для платформ:

```sh
wails build
```

Замечания
- В `frontend/index.html` в репозитории находится placeholder — замените его содержимым `scripts/tg_ui.html`.
- Wails позволяет экспортировать Go‑функции в фронтенд (open file dialog, локальное хранилище и т.д.).

Если хотите, я могу:
- Сгенерировать заготовку `desktop/wails` и скопировать UI (локально потребуется выполнить `wails init`).
- Добавить CI workflow для сборки на GitHub Actions.

``` 
EventyOn — Wails frontend scaffold notes

This folder is a lightweight scaffold for integrating the UI into a Wails project. Because the Wails CLI initializes several platform-specific files, this repo includes a placeholder frontend file and instructions to finish the setup locally.

Steps to finish and build a Wails app (multiplatform):

1. Install prerequisites
   - Go 1.18+ (you have 1.24)
   - Node.js + npm
   - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2. Initialize a Wails project inside `desktop/wails` (one-time):

   cd desktop/wails
   wails init -n eventyon-desktop -t vanilla

   This will create a `frontend/` and `backend/` structure. Replace `frontend/src/index.html` (or the template files) with the full contents of `../../scripts/tg_ui.html` (the file already present in the repo).

3. Run in development mode:

   cd desktop/wails
   wails dev

   The desktop window will open and load the embedded frontend. The frontend uses `http://localhost:8080` as the default API_BASE; ensure your backend is running.

4. Build for platforms:

   wails build

   For multi-platform CI builds, run `wails build` on each target OS runner (macOS, Windows, Linux) or use a cross-build environment.

Notes:
- The included `frontend/index.html` is a placeholder. Copy the contents of `scripts/tg_ui.html` into the Wails frontend index (or adapt the template) to reuse the UI.
- Wails allows exposing Go functions to the frontend via bindings; if you want tight integration (e.g., open file dialog, store files locally), we can add Go bindings later.
