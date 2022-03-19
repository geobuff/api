#!/bin/bash

if [ "$CIRCLE_BRANCH" = "main" ]; then
    echo "vpc_access_connector:" >> app.yaml
    echo "  name: projects/$PROD_GOOGLE_PROJECT_ID/locations/$PROD_REGION/connectors/$PROD_CONNECTOR_NAME" >> app.yaml
    echo "  egress_setting: all-traffic" >> app.yaml
fi
