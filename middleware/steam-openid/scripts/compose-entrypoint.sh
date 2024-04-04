#!/bin/sh
# export steam key file content to env var 
export STEAM_API_KEY=$(cat $STEAM_API_KEY_FILE)
npm run start