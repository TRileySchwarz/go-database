# go-database
This project contains a Golang database implementation that serves as an example to how someone might look at implementing docker / postgres / JWT etc. 

This API defaults to port 8080 for the app, and 5432 with the postgres. Ensure these ports are not in use.
The payloads are to be sent via a raw transaction, not URL encoded etc. 

From the project root

```$ docker-compose up```


Test can be run locally

``` $ go test ./...```





# API Documentation

### `POST /signup`
Endpoint to create an user row in postgres db.where `email` is an unique key in the database. The payload should have the following fields:

```json
{
  "email": "Chad@gmail.com",
  "password": "badPassword",
  "firstName": "Chad",
  "lastName": "Chillerton"
}
```

The response body returns a JWT on success that can be used for other endpoints:

```json
{
  "token": "some_jwt_token" 
}
```

### `POST /login` 

Endpoint to log an user in. 

The payload:

```json
{
  "email": "Chad@gmail.com",
  "password": "badPassword"
}
```

The response body returns a JWT on success that can be used for other endpoints:

```json
{
  "token": "some_jwt_token"
}
```

### `GET /users`
Endpoint to retrieve a json of all users. This endpoint requires a valid `x-authentication-token` header to be passed in with the request.

The response:
```json
{
  "users": [
    {
      "email": "Chad@gmail.com",
      "firstName": "Optimus",
      "lastName": "Prime"
    }
  ]
}
```

### `PUT /users`
Endpoint to update the current user `firstName` or `lastName` only. This endpoint requires a valid `x-authentication-token` header to be passed in and it should only update the user of the JWT being passed in. 

The payload:

```json
{
  "firstName": "NewFirstName",
  "lastName": "NewLastName"
}
```



### TODO Frontend
- Get the react app running on DO Droplet
- Switch the ports of react app to be port 81,
- Get the basic backend http Server handler to run on port 82(instead of 81 as currently)
- Update nginx router
- Ensure chadchillbro.com points to react app, not the basic backend anymore
- Get the API to frontend to make http call to port 82 backend
- Make https work via nginx
- Get go-database docker-compose working on the backend on port 82 and allow frontend access
- Combine frontend docker-compose file and go-database together
- Connect the front end in such a way that now only connections from within the server are allowed, ie API is not public
- Deny all port requests other than chadchillbro.com
- Create public api at api.chadchillbro.com   


### TODO Backend

- Implement rate limit for ip
- Add automated testing for docker-compose updates
- Build exec watcher to auto docker-compose build
- Verify Volume files are working as intended
- Write postman tests
- Provide command line args for port numbers


