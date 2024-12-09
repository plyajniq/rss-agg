# О проекте

Проект агрегатора RSS-лент

# Технологический стек

Бекенд:

    - Go 1.23
    - chi
    - goose 3.22
    - sqlc 1.27

Инфраструктура:

    - Linux, Ubuntu
    - Docker
    - Docker-Compose
    - Git Actions
    - Nginx
    - Swagger

База данных:

    - Postgres

Фронтенд:

    - HTML
    - CSS

# Автор

Никита Сергеевич Федяев

Telegram: [@nsfed](https://t.me/nsfed)

Репозиторий: [GitHub](git@github.com:plyajniq/rss-agg.git)

Проект доступен по адресу: [rss-aggregator.ru](https://rss-aggregator.ru/)

# Функционал

Проект работает в двух режимах

1. Веб-интерфейс для отображения RSS-лент

2. API с расширенным функционалом. [Документация эндпойнтов API.](http://rss-aggregator.ru/swagger/index.html)
    - Регистрация нового пользователя
    - Авторизация по токену
    - Получение информации о пользователе
    - Добавления новых лент
    - Подписка на ленты
    - Получение постов лент, на которые подписан пользователь
