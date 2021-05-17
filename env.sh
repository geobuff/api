#!/bin/bash

if [ "$CIRCLE_BRANCH" = "prod-pipeline" ]; then
    echo SITE_URL=$DEV_SITE_URL >> .env
    echo CONNECTION_STRING=$DEV_CONNECTION_STRING >> .env
    echo AUTH_ISSUER=$DEV_AUTH_ISSUER >> .env
    echo AUTH_SIGNING_KEY=$DEV_AUTH_SIGNING_KEY >> .env
    echo CORS_ORIGINS=$DEV_CORS_ORIGINS >> .env
elif [ "$CIRCLE_BRANCH" = "main" ]; then
    echo SITE_URL=$PROD_SITE_URL >> .env
    echo CONNECTION_STRING=$PROD_CONNECTION_STRING >> .env
    echo AUTH_ISSUER=$PROD_AUTH_ISSUER >> .env
    echo AUTH_SIGNING_KEY=$PROD_AUTH_SIGNING_KEY >> .env
    echo CORS_ORIGINS=$PROD_CORS_ORIGINS >> .env
fi

echo CORS_METHODS=$CORS_METHODS >> .env
echo CORS_HEADERS=$CORS_HEADERS >> .env
echo SENDGRID_API_KEY=$SENDGRID_API_KEY >> .env