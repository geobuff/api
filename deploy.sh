#!/usr/bin/env bash

echo "SITE_URL=$SITE_URL\n" >> .env
echo "CONNECTION_STRING=$CONNECTION_STRING\n" >> .env
echo "AUTH_SIGNING_KEY=$AUTH_SIGNING_KEY\n" >> .env
echo "AUTH_ISSUER=$AUTH_ISSUER\n" >> .env
echo "CORS_ORIGINS=$CORS_ORIGINS\n" >> .env
echo "CORS_METHODS=$CORS_METHODS\n" >> .env

echo $GCLOUD_SERVICE_KEY | gcloud auth activate-service-account --key-file=-
gcloud --quiet config set project ${GOOGLE_PROJECT_ID}
gcloud app deploy