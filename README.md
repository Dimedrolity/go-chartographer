# Тестовое задание на стажировку в Контур

## Задача

Реализовать сервис по работе с изображениями формата BMP, HTTP API, покрыть тестами. [Ссылка на исходное задание](README_task.md). Дедлайн по сдачи решения был 20 Марта 2022 23:59 Мск. 

## Сложность задания

Есть ограничение по оперативной памяти - `2 Гбайт`. Было реализовано хранение изображений по частям фиксированного размера (тайлам), чтобы не помещать полное изображение в оперативную память.

## Что хорошего в решении

Приложение разбито на *слои*: API, сервисы (реализация логики и бизнес логики), репозитории (хранилища). Благодаря разбиению, каждый слой, пакет, структуру можно разрабатывать и тестировать независимо.

Написаны модульные тесты для каждого пакета, итоговое покрытие `~80%`.

Есть .run, Makefile, Dockerfile для простоты запуска.

Старался следовать [go-project-layout](https://github.com/golang-standards/project-layout), директории cmd, internal, pkg.

Зависимости хранятся в репозитории, директория vendor.

## Как можно улучшить

### Добавление персистентного хранилища для данных изображений

При повторном запуске проекта, данные о предыдущих созданных изображениях будут потеряны, так как они хранятся в оперативной памяти. Имеется в виду данные изображений, содержащие Id изображений, размер и т. д., а не сами изображения формата BMP.

Также при запуске через Docker, изображения хранятся в контейнере на диске и не связаны с хостом с помощью `bind mount`.

### Поддержка других форматов изображений

Реализован формат BMP и цветовая модель RGB 24-бит. Теоретически можно добавить поддержку PNG, JPEG, и т. д., и цветовых моделей RGBA32, Grayscale, и т. д.

### Уменьшение объема Docker образа

Добавление [multi-stage build](https://docs.docker.com/develop/develop-images/multistage-build/) в Dockerfile.

### Добавление docker-compose

Для того, чтобы зафиксировать аргументы команд создания образа и запуска контейнера.

### Реализация TODO-комментариев в коде

Что-то есть, но не критично.

## Разное

Мой первый проект на Go, изучаю язык с февраля.

При решении начал с MindMap и почти сразу же создал Trello для ведения задач. 

Использовал Git (must have). 

На момент выполнения тестового был слабо знаком с Docker, это мой первый Dockerfile.