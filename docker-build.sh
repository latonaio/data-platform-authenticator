#!/bin/sh

PUSH=$1
DATE="$(date "+%Y%m%d%H%M")"
REPOSITORY_PREFIX="latonaio"
SERVICE_NAME="data-platform-authenticator"

PRIVATE_KEY="$(cat private.pem)"
PUBLIC_KEY="$(cat public.pem)"

# todo kubernetes の secret を検討した後必要であれが kubernetes secret を使う
# build時に環境変数PRIVATE_KEYに生成したprivate.keyの中身をセット
docker build --build-arg PRIVATE_KEY="${PRIVATE_KEY}" --build-arg PUBLIC_KEY="${PUBLIC_KEY}" -t ${SERVICE_NAME}:"${DATE}" -f Dockerfile .
# DOCKER_BUILDKIT=1 docker build --progress=plain -t ${SERVICE_NAME}:"${DATE}" .

# tagging
docker tag ${SERVICE_NAME}:"${DATE}" ${SERVICE_NAME}:latest
docker tag ${SERVICE_NAME}:"${DATE}" ${REPOSITORY_PREFIX}/${SERVICE_NAME}:"${DATE}"
docker tag ${REPOSITORY_PREFIX}/${SERVICE_NAME}:"${DATE}" ${REPOSITORY_PREFIX}/${SERVICE_NAME}:latest

if [[ $PUSH == "push" ]]; then
    docker push ${REPOSITORY_PREFIX}/${SERVICE_NAME}:"${DATE}"
    docker push ${REPOSITORY_PREFIX}/${SERVICE_NAME}:latest
fi
