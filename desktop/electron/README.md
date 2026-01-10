```markdown
# EventyOn — Electron оболочка (локальное тестирование)

В этой папке находится минимальный Electron‑враппер для локальной проверки интерфейса.

Быстрый старт (локально):

1. Подготовьте фронтенд
   - Для удобства можно скопировать содержимое `../../scripts/tg_ui.html` в `desktop/electron/public/index.html`.
   - Либо отредактируйте `desktop/electron/public/index.html`, вставив нужные скрипты и стили.

2. Установите зависимости и запустите:

```bash
cd desktop/electron
npm install
npm start
```

Откроется окно с интерфейсом. По умолчанию UI обращается к `http://localhost:8080` — убедитесь, что сервер запущен.

Сборка для релиза

Для production‑сборок используйте `electron-builder` или `electron-packager`. Текущий scaffold ориентирован на локальную разработку; при необходимости могу добавить конфиг для `electron-builder` и CI workflow.

Заметки

- Preload‑скрипт предоставляет небольшое безопасное API; все сетевые запросы выполняет рендерер (UI).
- Для продакшна проверьте настройки безопасности (CSP), запрет удалённого контента и подпись приложений на macOS.

``` 
# EventyOn — Electron desktop client (local testing)

This folder contains a minimal Electron wrapper so you can test the UI as a desktop app.

Quick start (local testing):

1. Copy the renderer UI into this folder
   - For convenience, you can copy the contents of `../../scripts/tg_ui.html` into `desktop/electron/public/index.html`.
   - Alternatively, edit `desktop/electron/public/index.html` and include the script and styles from `scripts/tg_ui.html`.

2. Install and run

```bash
cd desktop/electron
npm install
npm start
```

The window will open and load the UI. The UI uses `http://localhost:8080` as the default backend; make sure the server is running.

Builds

To produce production builds for specific platforms, use a packer such as `electron-builder` or `electron-packager`. This minimal scaffold is geared for local testing; I can add a full electron-builder config and CI workflow if you want cross-platform production artifacts.

Notes

- The preload script exposes a very small safe API; all network requests are done by the renderer (the UI).
- For production, review security options (CSP, disable remote content, signing on macOS, etc.).
