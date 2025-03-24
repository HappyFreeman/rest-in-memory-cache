
# 🚀 REST In-Memory Cache Service

> REST-сервис на Go с Fiber, реализующий CRUD-операции для хранения задач в памяти.

---

## 📖 **Описание**
Проект — это REST API на Go, использующий Fiber для создания кэша задач в памяти. Поддерживаются базовые CRUD-операции:

✅ Создание задачи  
✅ Получение задачи по ID  
✅ Получение списка всех задач  
✅ Обновление задачи  
✅ Удаление задачи

---

## 🏗️ **Структура проекта**
```
├── cmd
│   └── main.go         // Точка входа
├── docs
│   └── swagger
│       └── swagger.yaml    // Swagger-документация
├── internal
│   ├── api             // API-маршруты
│   │   └── middleware  // Middleware для обработки ошибок
│   ├── dto             // DTO-объекты для формирования ответов
│   ├── config          // Конфигурационные переменные
│   ├── repo            // Логика работы с памятью (репозиторий)
│   │   └── mock        // Мок-репозиторий
│   └── service         // Бизнес-логика
├── pkg
│   ├── logger          // Кастомный логгер
│   └── validator       // Валидаторы данных
├── go.mod              // Модуль Go
└── README.md
```

---

## 🛠️ **Стек технологий**
- Go
- Fiber
- Swagger
- zap (логирование)
- pkg/errors (обработка ошибок)

---

## 🚀 **Установка**
1. Склонируй репозиторий:
```bash
git clone https://github.com/HappyFreeman/rest-in-memory-cache.git
```

2. Перейди в папку проекта:
```bash
cd rest-in-memory-cache
```

3. Установи зависимости:
```bash
go mod tidy
```

---

## 🌐 **Запуск проекта**
Выполни команду для запуска:
```bash
go run cmd/main.go
```

---

## 📖 **Swagger-документация**
1. Сгенерируй Swagger-документацию:
```bash
swag init
```

2. Подключи Swagger в `main.go`:
```go
import "github.com/gofiber/swagger"

app.Get("/swagger/*", swagger.HandlerDefault)
```

3. Открой Swagger в браузере:  
   👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 📌 **API Эндпоинты**

### ✅ **Создать задачу**
```http
POST /tasks
```
**Request:**
```json
{
  "title": "Learn Go",
  "description": "Understand the basics of Go"
}
```
**Response:**
```json
{
  "status": "success",
  "data": {
    "task_id": 1
  }
}
```

---

### ✅ **Получить все задачи**
```http
GET /tasks
```
**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "title": "Learn Go",
      "description": "Understand the basics of Go"
    }
  ]
}
```

---

### ✅ **Получить задачу по ID**
```http
GET /tasks/{id}
```
**Response:**
```json
{
  "status": "success",
  "data": {
    "title": "Learn Go",
    "description": "Understand the basics of Go"
  }
}
```

---

### ✅ **Обновить задачу**
```http
PUT /tasks/{id}
```
**Request:**
```json
{
  "title": "Learn Go Advanced",
  "description": "Understand concurrency in Go"
}
```
**Response:**
```json
{
  "status": "success",
  "data": {
    "title": "Learn Go Advanced",
    "description": "Understand concurrency in Go"
  }
}
```

---

### ✅ **Удалить задачу**
```http
DELETE /tasks/{id}
```
**Response:**
```json
{
  "status": "success"
}
```

---

## 🚨 **Ошибки**
| Статус | Описание |
|--------|----------|
| `400`  | Неверный запрос |
| `404`  | Задача не найдена |
| `500`  | Внутренняя ошибка сервера |

---

## 🏆 **TODO**
- [ ] Добавить авторизацию
- [ ] Добавить тесты

---

## ❤️ **Автор**
[HappyFreeman](https://github.com/HappyFreeman)

---

🔥 **Готово к запуску!** 🚀
