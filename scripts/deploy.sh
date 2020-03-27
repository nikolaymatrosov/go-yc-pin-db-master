#!/bin/bash

if [[ -f .env ]]
then
  set -a
  . .env
  set +a
fi

if [[ ! -e "build" ]]; then
    mkdir "build"
else
    rm -rf "build"
    mkdir "build"
fi

cp *.go ./build
cp config.yaml ./build
cp go.mod ./build
cp go.sum ./build
rm build.zip || echo '';
(
    cd build;
    zip -r9 ../build.zip .
)

s3cmd put ./build.zip s3://$DEPLOY_BUCKET/build.zip \
  --access_key=$AWS_ACCESS_KEY_ID \
  --secret_key=$AWS_SECRET_ACCESS_KEY \
  --region=ru-central1 \
  --host=storage.yandexcloud.net \
  --host-bucket=\%\(bucket\)s.storage.yandexcloud.net

yc serverless function version create \
  --function-name=pin-db-master \
  --runtime golang114 \
  --entrypoint pin.PinHandler \
  --memory 128m \
  --execution-timeout 30s \
  --package-bucket-name $DEPLOY_BUCKET \
  --package-object-name build.zip\
  --service-account-id $SERVICE_ACCOUNT_ID \
  --environment DB_TYPE=$DB_TYPE,CLUSTER_ID=$CLUSTER_ID,TARGET_AZ=$TARGET_AZ