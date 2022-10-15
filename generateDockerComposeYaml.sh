#!/bin/bash

privateKey=$(cat private.pem)
mysqlUsername=$(jq -r .mysql.username env.json)
mysqlUserPassword=$(jq -r .mysql.password env.json)

replacedKey=$(awk 'BEGIN{ ORS = " " }{ print $0 }' private.pem)
export PRIVATE_KEY=$replacedKey
export MYSQL_USER_NAME=$mysqlUsername
export MYSQL_USER_PASSWORD=$mysqlUserPassword
rm -rf docker-compose.yaml; envsubst <"docker-compose-template.yaml"> "docker-compose.yaml";
