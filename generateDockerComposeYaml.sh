#!/bin/bash

privateKey=$(cat private.pem)
publicKey=$(cat public.pem)
mysqlUsername=$(jq -r .mysql.username env.json)
mysqlUserPassword=$(jq -r .mysql.password env.json)
authenticatorPort=$(jq -r .authenticator.port env.json)

replacedPrivateKey=$(awk 'BEGIN{ ORS = "\\n" }{ print $0 }' private.pem)
replacedPublicKey=$(awk 'BEGIN{ ORS = "\\n" }{ print $0 }' public.pem)
export PRIVATE_KEY=$replacedPrivateKey
export PUBLIC_KEY=$replacedPublicKey
export MYSQL_USER_NAME=$mysqlUsername
export MYSQL_USER_PASSWORD=$mysqlUserPassword
export AUTHENTICATOR_PORT=$authenticatorPort
rm -rf docker-compose.yaml; envsubst <"docker-compose-template.yaml"> "docker-compose.yaml";
