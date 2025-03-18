#!/bin/bash

HOST="localhost"      # Имя хоста
PORT="5432"           # Порт
DB_NAME="project-sem-1"  # Имя вашей базы данных
USER="validator"          # Имя пользователя
PASSWORD="val1dat0r"      # Пароль пользователя
SQLQUERY="CREATE TABLE IF NOT EXISTS prices (id INTEGER PRIMARY KEY, name varchar(30), category varchar(30), price DECIMAL(10, 2), create_date timestamp);"

export PGPASSWORD="$PASSWORD"

go build ./main.go
sudo apt install postgresql postgresql-client
psql -h $HOST -p $PORT -U $USER -d $DB_NAME  -c "$SQLQUERY"

unset PGPASSWORD