
You could easily run this project on local environment by using **docker-compose**. Below is the step to deploy the application:

If you on Linux, you could deploy the application by running the `init` bash scripts on deployment folder
```
$ ./deployments/init.sh
```

Else, you could manually run all service needed for deployment 

1. Deploy Certificate Manager Service
```
$ cd cert-manager && docker-compose up
```

2. After cert-manager is up, run the API service
```
$ cd ../deployments && docker-compose-up
```

3. After all service is up and running, you could view the API documentation on <span style="padding: 2px 4px; color: rgb(175, 78, 93); background-color: rgb(30, 15, 17)">localhost:3000/v1/docs/api</span>  or the project codebase documentation on <span style="padding: 2px 4px; color: rgb(175, 78, 93); background-color: rgb(30, 15, 17)">localhost:3000/v1/docs/code</span>

5. Postman collection for testing are available on **TBA**
