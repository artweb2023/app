# Тестовое задание

## Описание

Это веб-приложение позволяет управлять базой данных клиентов. Оно включает в себя функции аутентификации, отображения панели инструментов и API для входа, выхода и изменения статуса клиента.

1. Установите Go версия 1.20.2 на вашем компьютере, если он еще не установлен. 
2. Клонируйте репозиторий
3. Установите зависимостии go перейдя в проект:
   - go mod init <название проекта>
   - go mod tidy
4. Установите MySql
   - создайте базу данных : CREATE DATABASE "ваше название" CHARACTER SET = utf8mb4;
   - создайте пользователя: CREATE USER 'юзернайм'@'localhost' IDENTIFIED BY 'пароль';
   - установите права пользователя: GRANT ALL PRIVILEGES ON название базы.* TO 'юзернайм'@'localhost';
   - Импортируйте таблицы: mysql -u user -p password database < /пусть к папке/data/sql/schema.sql
   - Импортируйте синтезированные данные: mysql -u user -p -D database < /пусть к папке/data/sql/data.sql
5. В папке cmd/client в файле main.go необходимо прописать данные к базе данных : 
"UserName:Password@(localhost:3306)/Database Name?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true"
6. Для запуска проекта используйте команду: go run cmd/client/main.go cmd/client/handlers.go