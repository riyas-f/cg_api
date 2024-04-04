#!/bin/sh
# export steam key file content to env var 
export SMTP_CONFIG_PASSWORD=$(cat $SMTP_PASSWORD_FILE)
npm run start