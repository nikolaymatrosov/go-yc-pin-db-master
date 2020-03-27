#!/bin/sh

if [[ -f .env ]]
then
    set -a
    . .env
    set +a
fi

yc serverless function create --name=pin-db-master

yc serverless trigger create timer \
    --name=pin-db-master \
    --cron-expression="$CRON" \
    --invoke-function-name=pin-db-master \
    --invoke-function-tag="\$latest" \
    --invoke-function-service-account-id=$SERVICE_ACCOUNT_ID

