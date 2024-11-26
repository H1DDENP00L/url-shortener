# Go URL Shortener
URL Shortener — это простой сервис для сокращения длинных URL-адресов. Проект использует Go для бэкенда, PostgreSQL для хранения данных и предоставляет API для генерации коротких ссылок.

## Описание
Программа позволяет пользователям отправлять POST-запросы с длинным URL и получать короткую ссылку. <br>
Когда пользователь посещает короткую ссылку, происходит редирект на оригинальный URL.

## Технологии
- **Go** — основной ЯП
- **PostgreSQL** — бд для хранения ссылок
- **Regexp** — для валидации URL

## Установка и запуск
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/H1DDENP00L/url-shortener.git
   cd url-shortener

2. Убедитесь, что у вас установлены следующие компоненты:
- Go (для компиляции и запуска приложения).
- PostgreSQL (для работы с базой данных).
- Postman для удобной работы с запросами/ответами

3. Создайте базу данных и таблицу в PostgreSQL. Для этого выполните SQL-запросы:
    ```sql
    CREATE DATABASE urlShortener;
    
    CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url TEXT NOT NULL);
    ```

4. Настройте строку подключения к базе данных в файле main.go. Найдите строку:
    ```go
    connStr := "postgres://postgres:pass@localhost:5432/urlShortener?sslmode=disable"
    ```

5. Установите необходимые зависимости:
   ```bash
   go mod tidy
   ```

6. Запустите сервер:
   ```bash
   go run main.go
   ```
   
7. Программа будет доступна по адресу: http://localhost:8080.

8. Пример запроса:
    ```json
    {
        "url": "https://tech.wildberries.ru/courses/golang/application"
    }
    ```
9. Пример ответа:
    ```json
    {
        "short_url": "http://localhost:8080/hDbRq7"
    }
    ```

