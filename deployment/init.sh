#!/bin/bash
set -e

health_check() {
    # Endpoint for health check
    local ENDPOINT=$1

    while true; do
        # Send a request to the endpoint
        local RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $ENDPOINT)

        # Check if the endpoint responded with HTTP 200
        if [ $RESPONSE -eq 200 ]; then
            echo "Endpoint responded successfully. Exiting."
            exit 0
        else
            echo "Endpoint check failed with HTTP status: $RESPONSE. Retrying in 5 seconds."
            sleep 5
        fi
    done
}

# get the directory where the script is located
SOURCE=${BASH_SOURCE[0]}
while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )
  SOURCE=$(readlink "$SOURCE")
  [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done

DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )


echo Starting certificate manager service. Please wait...

# run cert manager in the background
cd ../cert-manager && docker-compose up -d

# cert manager health check
health_check "https://localhost:5500/health"

echo Certificate manager is up and running on localhost:5500

echo Starting the api. Please wait...
health_check "https://localhost:3000/health"

echo api is up and running on localhost:3000



