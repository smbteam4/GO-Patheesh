## Robot-Apocalypse
A docker based application.

**Technologies Used**
- Golang
- MongoDB
- Mongo Express
- Docker
- Docker compose
- Swagger

**FrameWork**
- Go Fiber


## How to setup the application

    docker-compose up -d

after the containers are successfully up, hit [localhost:8084](localhost:8084) in any of the web browser.
It will open mongo express. Create a new database named as **robot-apocalypse** 

#### Serve swagger
swagger documentation is added with this application. You can run the following command to check the request and response types

    swagger serve -F redoc swagger.yaml 


## Sample

**Add new survivor**

    curl --request POST \
    --url http://localhost:8080/api/v1/survivors \
    --header 'Content-Type: application/json' \
    --data '{
        "id" : "srv1",
        "name":"survivor1",
        "age":16,
        "location" : {
            "latitude" :152.024,
            "longitude" : 524.14
        }      
    }'

**Update Survivors**

    curl --request PUT \
    --url http://localhost:8080/api/v1/survivors \
    --header 'Content-Type: application/json' \
    --data '{
        "id" : "sdn1231264",
        "name":"survivor1",
        "age":16,
        "location" : {
            "latitude" :100,
            "longitude" : 100
        }      
    }'

**Mark as infected**

    curl --request PUT \
    --url http://localhost:8080/api/v1/survivors/infected \
    --header 'Content-Type: application/json' \
    --data '{
        "id" : "srv1",
        "reported_by" :"srv2"
    }'

**Infected percentage**

    curl --request GET \
    --url http://localhost:8080/api/v1/report/percentage \
    --header 'Content-Type: application/json'
   

**List of infected**

    curl --request GET \
    --url http://localhost:8080/api/v1/report/infected \
    --header 'Content-Type: application/json'

**List of Non-infected**

    curl --request GET \
    --url http://localhost:8080/api/v1/report/non-infected \
    --header 'Content-Type: application/json'

**Load robots list to db**

    curl --request POST \
    --url http://localhost:8080/api/v1/robots/load \
    --header 'Content-Type: application/json' 
    
**Load robots list**

    curl --request GET \
    --url http://localhost:8080/api/v1/robots/list \
    --header 'Content-Type: application/json' 
    