# MKV-stage-server
A REST api written in Golang for tracking songs and setlists for a music organization or church.
The api allows users to create songs, setlists and roles to schedule rehearsals and performances. 
Basic authentication is done via Json Webtokens and can be retrieved from the authorization server.

## Features
 - Authorization API using NodeJS and Express to serve JWT
   - Postgres database for user management
   - Redis database for storing refresh tokens
 - Golang Backend API using [Gin](https://gin-gonic.com/) and [GORM](https://gorm.io/) that accepts JWT from the authorization server
   - MySQL database 

## Usage

Docker compose is needed to build and connect to the applications and databases
```
docker-compose -f docker-compose.yml up --build
```
## Endpoints

### Authorization
Authorization endpoints are located by default at ```PORT=9080```

| Endpoint      | Type | Description               | Body Fields     | Query      | JWT |
|---------------|------|---------------------------|-----------------|------------|-----|
| /users        | POST | Create a new user         | email, password |            | No  |
| /users        | GET  | Retrieve all users        |                 | ids,emails | No  |
| /users/me     | GET  | Retrieve user information |                 |            | Yes |
| /users/logout | GET  | Logs the user out         |                 |            | Yes |
| /users/:id    | GET  | Retrieve user with id     |                 |            | No  |

| Endpoint        | Type | Description            | Body Fields               | Query | JWT |
|-----------------|------|------------------------|---------------------------|-------|-----|
| /tokens/login   | POST | Signs the user in      | senderId, email, password |       | No  |
| /tokens/refresh | POST | Renews an access token | refresh                   |       | No  |

### Backend
Backend endpoints are located by default at ```PORT=8008```

| Endpoint      | Type | Description         | Body Fields                          | Query | JWT |
|---------------|------|---------------------|--------------------------------------|-------|-----|
| /users/create | POST | Creates a new user  | first_name, last_name, profile_color |       | Yes |
| /users        | GET  | Retrieves all users |                                      |       | No  |

| Endpoint         | Type   | Description                | Body Fields                          | Query | JWT |
|------------------|--------|----------------------------|--------------------------------------|-------|-----|
| /users/me        | GET    | Retrieves user information |                                      |       | Yes |
| /users/me/update | PUT    | Updates user information   | first_name, last_name, profile_color |       | Yes |
| /users/me/delete | DELETE | Remove user*               | id                                   |       | Yes |

\* Can remove other users if permission level is >= 3

Other endpoints here...
