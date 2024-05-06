# #!/bin/bash

# health_check() {
#     # Endpoint for health check
#     local ENDPOINT=$1

#     while true; do
#         # Send a request to the endpoint
#         local RESPONSE=$(curl -k -s -o /dev/null -w "%{http_code}" $ENDPOINT)

#         # Check if the endpoint responded with HTTP 200
#         if [ $RESPONSE -eq 200 ]; then
#             echo "Endpoint responded successfully. Exiting."
#             return
#         else
#             echo "Endpoint check failed with HTTP status: $RESPONSE. Retrying in 5 seconds."
#             sleep 5
#         fi
#     done
# }

# # retrieve GCP instance Public IP
# ZONE="asia-southeast2-a"

# # get the directory where the script is located
# SOURCE=${BASH_SOURCE[0]}
# while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
#   DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )
#   SOURCE=$(readlink "$SOURCE")
#   [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
# done

# DIR=$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )

# # set up secret file
# PASSWORD_HASH_SECRET_KEY_FILE=$DIR/../api/secrets/password_hash_key.txt
# JWT_SECRET_KEY_FILE=$DIR/../api/secrets/jwt_secret_key.txt
# DB_ACCOUNT_SECRET_FILE=$DIR/db/secrets/db_account_password.txt
# DB_AUTH_SECRET_FILE=$DIR/db/secrets/db_auth_password.txt
# DB_GAMES_SECRET_FILE=$DIR/db/secrets/db_games_password.txt
# DB_SESSION_SECRET_FILE=$DIR/db/secrets/db_session_password.txt
# STEAM_API_KEY_FILE=$DIR/../middleware/steam-openid/secrets/steam_api_key.txt
# SMTP_PASSWORD_FILE=$DIR/../middleware/mail/secrets/smtp_password.txt

# mkdir -p $DIR/../api/secrets
# mkdir -p $DIR/db/secrets
# mkdir -p $DIR/../middleware/mail/secrets
# mkdir -p $DIR/../middleware/steam-openid/secrets
# mkdir -p $DIR/../middleware/mail/secrets

# set -e

# INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
# export HOST=$(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --format='get(networkInterfaces[0].accessConfigs[0].natIP)')

# # don't change when the file already exists
# if ! [ -f $DB_ACCOUNT_SECRET_FILE ]; then 
#     echo "$(openssl rand -base64 128)" | tr -d '\n'> $DB_ACCOUNT_SECRET_FILE
# fi 

# if ! [ -f $DB_AUTH_SECRET_FILE ]; then 
#     echo "$(openssl rand -base64 128)" | tr -d '\n'> $DB_AUTH_SECRET_FILE
# fi 

# if ! [ -f $DB_GAMES_SECRET_FILE ]; then 
#     echo "$(openssl rand -base64 128)" | tr -d '\n'> $DB_GAMES_SECRET_FILE
# fi 

# if ! [ -f $DB_SESSION_SECRET_FILE ]; then 
#     echo "$(openssl rand -base64 128)" | tr -d '\n'> $DB_SESSION_SECRET_FILE
# fi 

# echo "$(openssl rand -base64 128)" | tr -d '\n'> $PASSWORD_HASH_SECRET_KEY_FILE
# echo "$(openssl rand -base64 128)" | tr -d '\n'> $JWT_SECRET_KEY_FILE


# # access secret manager
# gcloud secrets versions access latest --secret=STEAM_API_KEY > $STEAM_API_KEY_FILE
# gcloud secrets versions access latest --secret=SMTP_CONFIG_PASSWORD > $SMTP_PASSWORD_FILE

# echo $DIR

# GCP_CERT_FILE_BUCKET_URL=https://storage.googleapis.com/root-cert-bucket/my-root-cert.crt
# GCP_PRIVATE_KEY_SECRET_NAME=ROOT_CA_PRIVATE_KEY
# GCP_PRIVATE_KEY_PASSPHRASE_SECRET_NAME=ROOT_CA_KEY_PASSPHRASE

# curl ${GCP_CERT_FILE_BUCKET_URL} -o ${ROOT_CA_VOLUME}/root-ca.crt
# gcloud secrets versions access latest --secret="${GCP_PRIVATE_KEY_SECRET_NAME}" > ${ROOT_CA_VOLUME}/root-ca.key
# gcloud secrets versions access latest --secret="${GCP_PRIVATE_KEY_PASSPHRASE_SECRET_NAME}" > ${ROOT_CA_VOLUME}/passphrase

# echo Starting certificate manager service. Please wait...

# # run cert manager in the background
# cd $DIR/../cert-manager && docker-compose down &&  docker-compose build --no-cache && docker-compose up -d

# # cert manager health check
# health_check "https://localhost:5500/health"

# echo Certificate manager is up and running on localhost:5500

# echo Starting the api. Please wait...
# cd $DIR && docker-compose down && docker-compose build --no-cache && docker-compose up -d
# health_check "https://localhost:3000/health"
# echo api is up and running on localhost:3000
# echo success
# exit 0


# Default values
skipBuild=false
noCache=false
secure=false
deployment="local"
instanceHost=""

usage() {
    echo "Usage: $0 [-s|--skipBuild] [-n|--noCache] [-h|--help]"
    echo "Options:"
    echo "  -s, --skipBuild    Skip the build process"
    echo "  -n, --noCache      Do not use cache during build"
    echo "  -h, --help         Display this help message"
    exit 1
}



# Parse command line options
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
        -s|--skipBuild)
            skipBuild=true
            shift
            ;;
        -n|--noCache)
            noCache=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Error: Unknown option: $key"
            usage
            ;;
    esac
done

DIR=$(dirname "$0")

# Define paths
PASSWORD_HASH_SECRET_KEY_FILE="$DIR/../api/secrets/password_hash_key.txt"
JWT_SECRET_KEY_FILE="$DIR/../api/secrets/jwt_secret_key.txt"
DB_ACCOUNT_SECRET_FILE="$DIR/db/secrets/db_account_password.txt"
DB_AUTH_SECRET_FILE="$DIR/db/secrets/db_auth_password.txt"
DB_CALL_SECRET_FILE="$DIR/db/secrets/db_call_password.txt"
SMTP_PASSWORD_FILE="$DIR/../middleware/mail/secrets/smtp_password.txt"

# GCP Configurations
ZONE="asia-southeast2-a" # the zone of the GCP vm instances
GCP_CERT_FILE_BUCKET_URL="https://storage.googleapis.com/root-cert-bucket/my-root-cert.crt"
GCP_PRIVATE_KEY_SECRET_NAME="ROOT_CA_PRIVATE_KEY"
GCP_PRIVATE_KEY_PASSPHRASE_SECRET_NAME="ROOT_CA_KEY_PASSPHRASE"

echo "Deployment: $deployment"
echo "Use Cache: $noCache"

# Retrieve GCP instance Public IP
if [ "$deployment" = "GCP" ]; then
    INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" -s "http://metadata.google.internal/computeMetadata/v1/instance/name")
    INSTANCE_HOST=$(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --format='get(networkInterfaces[0].accessConfigs[0].natIP)')
elif [ -z "$INSTANCE_HOST" ]; then
    if [ -n "$instanceHost" ]; then
        INSTANCE_HOST=$instanceHost
    else
        read -p "Host environment variable is empty. Please input the host of the machine:" read
        if [ -z "$read" ]; then
            echo "Received empty string as host. Defaulting to 127.0.0.1"
            read="127.0.0.1"
        fi
        INSTANCE_HOST=$read
    fi
fi

if [ -n "$useHTTPS" ]; then
     SECURE=1
     SCHEME="https"
else
    SECURE=0
     SCHEME="http"
fi

echo "HOST: $INSTANCE_HOST"
echo "DB_VOLUME: $DB_VOLUME"
echo "ROOT_CA_VOLUME: $ROOT_CA_VOLUME"
echo "SECURE: $SECURE"

# Create directories if not exists
mkdir -p $DIR/../api/secrets
mkdir -p $DIR/db/secrets
mkdir -p $DIR/../middleware/mail/secrets
mkdir -p $DIR/../middleware/steam-openid/secrets
mkdir -p $DIR/../middleware/mail/secrets

# Generate secret files if not exists
if [ ! -f $DB_ACCOUNT_SECRET_FILE ]; then
    base64 /dev/urandom | head -c 128 > $DB_ACCOUNT_SECRET_FILE
fi

if [ ! -f $DB_AUTH_SECRET_FILE ]; then
    base64 /dev/urandom | head -c 128 > $DB_AUTH_SECRET_FILE
fi

if [ ! -f $DB_CALL_SECRET_FILE ]; then
    base64 /dev/urandom | head -c 128 > $DB_CALL_SECRET_FILE
fi

if [ "$skipBuild" = false ]; then
    base64 /dev/urandom | head -c 128 > $PASSWORD_HASH_SECRET_KEY_FILE
    base64 /dev/urandom | head -c 128 > $JWT_SECRET_KEY_FILE
fi

# Access Secret Manager
gcloud secrets versions access latest --secret=SMTP_CONFIG_PASSWORD > $SMTP_PASSWORD_FILE

curl $GCP_CERT_FILE_BUCKET_URL -o "$ROOT_CA_VOLUME/root-ca.crt"
gcloud secrets versions access latest --secret="$GCP_PRIVATE_KEY_SECRET_NAME" > "$ROOT_CA_VOLUME/root-ca.key"
gcloud secrets versions access latest --secret="$GCP_PRIVATE_KEY_PASSPHRASE_SECRET_NAME" | tr -d '\r\n' > "$ROOT_CA_VOLUME/passphrase"

echo "Starting certificate manager service. Please wait..."

# Run cert manager in the background
cd $DIR/../cert-manager
docker-compose down

if [ "$skipBuild" = false ]; then
    if [ "$noCache" = true ]; then
        docker-compose build --no-cache
    else
        docker-compose build
    fi
fi

docker-compose up -d

# Cert manager health check
health_check "https://localhost:5500/health"

echo "Certificate manager is up and running on localhost:5500"

echo "Starting the API. Please wait..."
cd $DIR

docker-compose down

if [ "$skipBuild" = false ]; then
    if [ "$noCache" = true ]; then
        docker-compose build --no-cache
    else
        docker-compose build
    fi
fi

docker-compose up -d

health_check "${SCHEME}://localhost:3000/health"
echo "API is up and running on localhost:3000"
echo "Success"
exit 0