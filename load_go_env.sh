#!/bin/bash

# Очистите файл перед началом
> /etc/default/go-backend.env

# Читаем переменные окружения из .env-vars и экспортируем их
while IFS='=' read -r key value; do
    echo "export $key=\"$value\"" >> /etc/default/go-backend.env
done < /home/2cubes/backend/.env-vars
