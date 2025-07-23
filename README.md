# Subscription Service

Микросервис для управления подписками пользователей.  
Поддерживает создание, обновление, удаление подписок, а также подсчёт суммарной стоимости за определённый период.

## **Функционал**
- **CRUD для подписок** (`/subscriptions`)
- **Подсчёт суммы подписок** (`/subscriptions/sum?user_id=...&service_name=...&start=MM-YYYY&end=MM-YYYY`)
- **Валидация UUID для user_id**
- **Поддержка Docker и docker-compose**

---

## **Стек технологий**
- **Go (Golang)** — основная логика
- **PostgreSQL** — база данных
- **Docker** и **docker-compose** — контейнеризация
- **gorilla/mux** — роутер
- **uuid** — для генерации и валидации UUID
- **sql.DB** — работа с БД

---

Примеры API
Создать подписку
```
POST /subscriptions
Content-Type: application/json

{
  "service_name": "Netflix",
  "price": 400,
  "user_id": "c53e7547-afd5-4ec5-ac97-c2e0fe229d82",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```
Получить подписку по ID
```
GET /subscriptions/18
```
Сумма подписок
```
GET /subscriptions/sum?user_id=c53e7547-afd5-4ec5-ac97-c2e0fe229d82&service_name=Netflix&start=07-2025&end=12-2025
```
##**Подробная документация изложена в файле swagger.yaml**##
