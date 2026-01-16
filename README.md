# 🗂 TaskManagerITC

**TaskManagerITC** — full‑stack приложение для управления задачами с backend на **Go** и frontend на **React**. Проект реализован с чистым разделением ответственности, поддерживает авторизацию через JWT, работу с БД и полностью готов к запуску в Docker.

## 📖 Описание проекта

**TaskManagerITC** — это учебно‑практический проект, демонстрирующий:

- построение REST API на Go
- работу с JWT‑авторизацией
- миграции базы данных
- современный frontend на React
- контейнеризацию через Docker и Docker Compose

Проект подходит как пример production‑подобной архитектуры для full‑stack приложений.

---

## 🧱 Структура проекта

```
TaskManagerITC/
├── backend/                     # Backend (Go)
│   ├── cmd/
│   │   └── api/                 # Точка входа в приложение
│   ├── internal/                # Внутренняя логика приложения
│   │   ├── handlers             # HTTP‑обработчики
│   │   ├── services             # Бизнес‑логика
│   │   ├── notificatios         # Настройка уведомлений
│   │   ├── database             # База данных
│   │   ├── repository           # Работа с БД
│   │   └── models               # Модели данных
│   ├── pkg/
│   │   └── jwt/                 # JWT‑утилиты (токены, middleware)
│   ├── migrations/              # SQL‑миграции
│   ├── .env                     # Переменные окружения
│   ├── Dockerfile               # Docker‑образ backend
│   ├── go.mod                   # Go‑зависимости
│   └── go.sum                   # Контрольные суммы
│
├── frontend/                    # Frontend (React)
│   ├── public/                  # Публичные файлы
│   ├── src/
│   │   ├── pages/               # Страницы приложения
│   │   ├── styles/              # CSS / стили
│   │   ├── App.js               # Корневой компонент
│   │   ├── index.js             # Точка входа
│   │   └── index.css            # Глобальные стили
│   ├── Dockerfile               # Docker‑образ frontend
│   ├── package.json             # JS‑зависимости
│   ├── package-lock.json
│   └── .env
│
├── docker-compose.yml            # Оркестрация сервисов
├── .gitignore
└── README.md
```

---

## 🚀 Функциональность

- 🔐 Авторизация и аутентификация (JWT)
- 👤 Работа с пользователями
- ✅ CRUD‑операции с задачами
- 🗃 Хранение данных в БД
- 🌐 REST API
- 🖥 SPA‑интерфейс на React
- 🐳 Полный запуск через Docker

---

## 🧠 Стек технологий

### 🔧 Backend

- **Go (Golang)**
- **net/http** или роутер (chi / gin)
- **JWT** — авторизация
- **SQL (PostgreSQL / MySQL)**
- **golang‑migrate** — миграции
- **Docker**

### 🎨 Frontend

- **React**
- **JavaScript (ES6+)**
- **SCSS**
- **CSS**
- **Fetch API / Axios**

### 🐳 DevOps

- **Docker**
- **Docker Compose**

---

## ⚙️ Переменные окружения (Backend)

Пример `.env`:

```env
APP_PORT=8080
TELEGRAM_BOT_TOKEN=
JWT_SECRET=
JWT_TTL=24h
NAME_OF_DATABASE=backend/internal/model/database/projects_db.db
DATABASE=sqlite3
```

---

## 📦 Загрузка пользователей в БД

CSV импортируется локально и записывается в таблицу `users`:

```bash
USERS_CSV_PATH=/path/to/users.csv python3 backend/scripts/seed_users.py
```

---

## 🐳 Запуск через Docker

```bash
docker-compose up --build
```

После запуска:

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`

---

## ▶️ Локальный запуск

### Backend

```bash
cd backend
go mod download
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm start
```

---

## 📡 API (пример)

```http
GET   /health
POST   /auth/telegram
GET    /get_users
```

Все защищённые эндпоинты требуют JWT‑токен.

---

## 🏗 Архитектура

Проект следует принципам:

- Clean Architecture
- Dependency Injection
- Разделение слоёв (handler → service → repository)
- Backend и frontend полностью изолированы

---

## 📌 Статус проекта

🟢 В разработке / учебный проект айти сообщества

---
