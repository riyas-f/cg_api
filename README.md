# cloud-gaming-api
API for cloud gaming service

## Installation
1. Install Docker and Docker Compose
2. Set the neccessary environment variables
```
export INSTANCE_HOST=127.0.0.1
export ROOT_CA_VOLUME=/some/path/to/folder/with/root/ca
export DB_VOLUME=/some/path/to/db/volume/folder
export SCHEME=http
export SECURE=0
```
4. Build and Run Certificate Manager
```
$ cd cert-manager && docker-compose build && docker-compose up -d
```
3. Build and run the API
```
$ cd ../deployment && docker-compose build && docker-compose up -d
```
4. API will be available on port 3000
