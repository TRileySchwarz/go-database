# go-database
This project contains a Golang database implementation that serves as an example to how someone might look at implementing docker / postgres / JWT etc. 



## API Specs

### `POST /signup`
Endpoint to create an user row in postgres db. The payload should have the following fields:

```json
{
  "email": "Chad@gmail.com",
  "password": "badPassword",
  "firstName": "Chad",
  "lastName": "Chillerton"
}
```

where `email` is an unique key in the database.

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





## TODO

- Add password hashing so they are not stored in plain text
- Implement rate limit for ip
- Add better error handling
- Add automated testing for docker-compose updates
- Build super rough UI / front end service in react
- Nginx reverse proxy 
- Forward DNS
- Build exec watcher to auto docker-compose build
- Password verification of difficulty, eg. longer than 10 characters
- Verify Volume files are working as intended
- Write postman tests

