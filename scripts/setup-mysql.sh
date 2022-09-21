#!/bin/bash
DB_USER_NAME="XXXXXXXX"
DB_USER_PASSWORD="XXXXXXXX"
SRC_DIR="/app/src"
CONF_PATH="$SRC_DIR/mysql.conf"
CONTAINER_NAME="sample-mysql"

# CONTAINER_ID=echo docker container ls --filter "name=$CONTAINER_NAME" --quiet
# echo $CONTAINER_ID

docker exec -it ${CONTAINER_NAME} sh -c "mysql --defaults-extra-file=${CONF_PATH} -t --show-warnings -e \"CREATE USER $DB_USER_NAME@localhost IDENTIFIED BY '$DB_USER_PASSWORD';\" "