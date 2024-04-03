#!/bin/sh
# export steam key file content to env var 
export DB_PASSWORD=$(cat $DB_PASSWORD_FILE)

# export related stuff
if [[ -n $HASH_SECRET_KEY_FILE]]; then 
    export HASH_SECRET_KEY=$(cat $HASH_SECRET_KEY_FILE)
elif [[ -n $JWT_SECRET_KEY_FILE]]; then 
    export JWT_SECRET_KEY=$(cat $HASH_SECRET_KEY_FILE)
fi

./app