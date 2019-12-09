#!/bin/sh

docker run --name mariadb -v ~/novabot/mariadb_storage:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=development -p 3306:3306 -d mariadb:latest

