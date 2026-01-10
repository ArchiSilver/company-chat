# Company Chat - Полнофункциональная коммуникационная платформа

[![React Native](https://img.shields.io/badge/React%20Native-0.83.1-blue.svg)](https://reactnative.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8.3-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)

Company Chat - это полнофункциональная коммуникационная платформа, включающая веб-приложение, мобильное приложение и бэкенд API. Проект разработан для корпоративного использования с современным интерфейсом и высокой производительностью.

## 📁 Структура проекта

```
company-chat/
├── cmd/                    # Исполняемые файлы
│   ├── app/               # Основное веб-приложение
│   └── migrate/           # Миграции базы данных
├── internal/              # Внутренняя логика (Go)
│   ├── config/           # Конфигурация
│   ├── domain/           # Бизнес-логика
│   ├── handlers/         # HTTP обработчики
│   └── repository/       # Работа с БД
├── migrations/           # SQL миграции
├── web/                  # Веб-интерфейс (React)
├── mobile/               # Мобильное приложение (React Native)
├── EventyOn/             # Отдельное мобильное приложение
├── desktop/              # Десктопные приложения (Electron)
├── scripts/              # Скрипты разработки
└── docs/                 # Документация
```

## 🚀 Быстрый старт

### Предварительные требования

- **Go 1.21+**
- **Node.js 20+**
- **PostgreSQL 15+**
- **Android Studio** (для мобильной разработки)
- **Xcode** (для iOS разработки)

### Установка и запуск

1. **Клонировать репозиторий:**
   ```bash
   git clone https://github.com/ArchiSilver/company-chat.git
   cd company-chat
   ```

2. **Настроить бэкенд:**
   ```bash
   # Установить зависимости
   go mod download

   # Запустить PostgreSQL (через Docker)
   make postgres

   # Запустить миграции
   make migrate-up

   # Запустить сервер
   make run
   ```

3. **Настроить веб-приложение:**
   ```bash
   cd web
   npm install
   npm run dev
   ```

4. **Настроить мобильное приложение:**
   ```bash
   cd mobile
   npm install
   npm start
   # В новом терминале:
   npm run android  # или npm run ios
   ```

## 📱 Мобильные приложения

### EventyOn Messenger

Современное мобильное приложение для мгновенного обмена сообщениями с элегантным интерфейсом.

#### Особенности:
- 🎨 **Современный дизайн** с фиолетовой цветовой схемой (#8A2BE2)
- 📱 **Кроссплатформенность** (iOS/Android)
- 🔐 **Регистрация** по телефону или email
- 💬 **Мессенджер** в стиле Telegram
- ⚡ **Быстрая навигация** между чатами
- 🌙 **Адаптивный интерфейс**

#### Технологии:
- **React Native 0.83.1**
- **TypeScript**
- **React Navigation**
- **React Native Paper**
- **AsyncStorage**

#### Запуск EventyOn:
```bash
cd EventyOn
npm install
npm start
# В новом терминале:
npx react-native run-android
```

### Company Chat Mobile

Полнофункциональное мобильное приложение для корпоративного общения.

#### Особенности:
- 👥 **Управление пользователями**
- 💼 **Корпоративные чаты**
- 📎 **Отправка файлов**
- 🔒 **Аутентификация**
- 📊 **Статистика сообщений**

## 🖥️ Веб-приложение

### Технологии:
- **React 18**
- **TypeScript**
- **Vite**
- **Tailwind CSS**
- **WebSocket** для real-time обновлений

### Запуск:
```bash
cd web
npm install
npm run dev
```

## 🗄️ Бэкенд API

### Технологии:
- **Go 1.21+**
- **Gin** (HTTP фреймворк)
- **PostgreSQL** (база данных)
- **Redis** (кэширование)
- **WebSocket** (real-time коммуникация)
- **JWT** (аутентификация)

### API Endpoints:

#### Аутентификация:
- `POST /api/auth/login` - Вход
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/refresh` - Обновление токена

#### Чаты:
- `GET /api/chats` - Список чатов
- `POST /api/chats` - Создать чат
- `GET /api/chats/:id/messages` - Сообщения чата
- `POST /api/chats/:id/messages` - Отправить сообщение

#### Пользователи:
- `GET /api/users` - Список пользователей
- `GET /api/users/:id` - Информация о пользователе
- `PUT /api/users/:id` - Обновить профиль

## 🐳 Docker

### Быстрый запуск с Docker:
```bash
# Запустить все сервисы
docker-compose up -d

# Посмотреть логи
docker-compose logs -f
```

### Сервисы:
- **app**: Go веб-сервер (порт 8080)
- **postgres**: PostgreSQL база данных (порт 5432)
- **redis**: Redis кэш (порт 6379)

## 📊 База данных

### Структура:
- **users**: Пользователи системы
- **chats**: Чаты и группы
- **messages**: Сообщения
- **uploads**: Загруженные файлы

### Миграции:
```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down
```

## 🧪 Тестирование

### Запуск тестов:
```bash
# Unit тесты
make test

# Integration тесты
make test-integration

# E2E тесты
make test-e2e
```

### Мобильное тестирование:
```bash
cd mobile
npm test
```

## 📦 Сборка и развертывание

### Сборка бэкенда:
```bash
make build
```

### Сборка мобильного приложения:
```bash
cd mobile
npm run build:android  # или build:ios
```

### Docker сборка:
```bash
docker build -t company-chat .
```

## 🔧 Разработка

### Скрипты:
- `make dev` - Запуск в режиме разработки
- `make dev-mobile` - Разработка мобильного приложения
- `scripts/dev-up.sh` - Полный запуск всех сервисов
- `scripts/preview.sh` - Предпросмотр веб-приложения

### Переменные среды:
```bash
# .env файл
DATABASE_URL=postgres://user:password@localhost:5432/companychat
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
```

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для вашей фичи (`git checkout -b feature/AmazingFeature`)
3. Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`)
4. Запушьте в ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект лицензирован под MIT License - смотрите файл [LICENSE](LICENSE) для деталей.

## 👥 Авторы

- **ArchiSilver** - *Основной разработчик*

## 🙏 Благодарности

- React Native Community
- Gin Web Framework
- PostgreSQL Community
- Все контрибьюторы проекта

## 📞 Контакты

- **Email**: archisilver@example.com
- **GitHub**: [@ArchiSilver](https://github.com/ArchiSilver)
- **LinkedIn**: [ArchiSilver](https://linkedin.com/in/archisilver)

---

⭐ **Если проект вам понравился, поставьте звезду!** ⭐
