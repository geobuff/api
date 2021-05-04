#!/bin/bash

echo SITE_URL=$SITE_URL >> .env
echo DATABASE_HOST=$DATABASE_HOST >> .env
echo DATABASE_USER=$DATABASE_USER >> .env
echo DATABASE_PASSWORD=$DATABASE_PASSWORD >> .env
echo DATABASE_NAME=$DATABASE_NAME >> .env
echo AUTH_SIGNING_KEY=$AUTH_SIGNING_KEY >> .env
echo AUTH_ISSUER=$AUTH_ISSUER >> .env
echo CORS_ORIGINS=$CORS_ORIGINS >> .env
echo CORS_METHODS=$CORS_METHODS >> .env

echo $GCLOUD_SERVICE_KEY | gcloud auth activate-service-account --key-file=-
gcloud --quiet config set project ${GOOGLE_PROJECT_ID}
gcloud app deploy