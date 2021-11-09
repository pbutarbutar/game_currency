# Go-JWT-Mysql-Restful-API

How To Run Docker

1. Please do build first
    --> docker-compose build
2. For to running All service do Docker Compose up
    --> docker-compose up -d

In the docker compose there is service for ready to installed
1. Go
2. MySQL

After 2 steps runing completed
Port Expose : 8083 and apps already can to using

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

