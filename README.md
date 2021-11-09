# Go-JWT-Mysql-Restful-API

How To build docker and run all service in docker compose.
1. Please do build first
    --> docker-compose build 
2. For to running All service do Docker Compose up
    --> docker-compose up -d
3. All container services it's going be running well and can to use, port expose 8083

For to run test unit we only change docker-compose.test.yml to docker-compose.yml

In the docker compose there is service is already success to install
1. Go
2. MySQL

All endpoint there is already to use.
How to use enpdoint please follow like below,
1. localhost:8083/login {POST}
    payloads in the body should json format
     {
        "email":"p1@gmail.com",
        "password":"password"
     }
2. localhost:8083/customers {POST}
    payloads in the body should json format
     {
        "name":"Parulian4",
        "author_id":1
     }
     and set authorization bearer a token from result login (Required)
3. localhost:8083/currencies {POST}
    payloads in the body should json format
     {
        "currency_from":1,
        "currency_to":2,
        "rate":29
     }
     and set authorization bearer a token from result login (Required)

4. localhost:8083/calculatecurrency {POST}
    payloads in the body should json format
     {
        "currency_from":1,
        "currency_to":2,
        "amount":580
    }
     and set authorization bearer a token from result login  (Required)      

Documentation also create in swagger, apiary.io

My Experience i have create API documentation in apiary a API GraphQL --> https://pintekid.docs.apiary.io/#

