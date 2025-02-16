# test-task

## Это репозиторий тестового задания авито

Чтобы запустить проект необходимо:

1. Склонировать репозиторий
2. Открыть директорию с склонированным репозиторием
3. docker-compose up --build чтобы начать запуск сервера
4. Спустя некоторое время сервер буде досутпен по порту 8080
5. Для остановки сервера можно использовать Ctrl+C

## Linter

В корневой директории проекат лежит файл `golangci.yml` который использовался во время проекта как конфиг линтера golangci

## Нагрузочное тестирование

Тестирование было произведено с помощью vegeta

Команда для запуска: `cat targets.txt | vegeta attack -rate=1000 -duration=5s | tee results.bin | vegeta report`

Пример файла targets.txt:

```
GET http://localhost:8080/api/info
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk4MjIzNDIsImlhdCI6MTczOTczNTk0MiwidXNlcm5hbWUiOiJhYWFhYSJ9.XTA-qrF_xFuvT-36qwgnusTk0rj_ORusq0Yh-Kg3eHE
```

Результаты тестирования:

![image_2025-02-16_23-03-10](https://github.com/user-attachments/assets/26aef090-4441-4ebe-ab84-e0a9e0ea0349)
