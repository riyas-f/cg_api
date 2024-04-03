#!/bin/sh
# export steam key file content to env var 
export DB_PASSWORD=$(cat $DB_PASSWORD_FILE)
./app