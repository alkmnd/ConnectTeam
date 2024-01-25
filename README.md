# ConnectTeam server app 

## Run 
1. First, install go and golang-migrate. For macOS:
``` bash
brew install golang-migrate
```
2. Add .env file in the project.
3. Run docker container. Examle:
``` bash
docker run --name=ct-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
```
4. Run migrations:
``` bash
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
```
5. Run app:
``` bash
go run cmd/main.go
```
## API Documentation

Note:

For all REST Api authenticated requests, ensure to include the Authorization header with the value Bearer 
your_access_token. This token is obtained through the user authentication process described in the
first section.
All requests require a valid authorization token in the header.

### 1. Authentication 
#### 1.1 User Sign-Up 

Method: `POST`

Endpoint: `/auth/sign-up`

Description: Registrates user

Request Parameters:  

*`email` (string, required): User email  

*`phone_number` (string, required): User phone number 

*`first_name` (string, required)

*`second_name` (string, required)

*`password` (string required)

Response: 
* `id` (string): Signed-up user id

Example body: 
``` bash
{
    "email": "dkhfлg@gmail.com",
    "phone_number": "89912818155",
    "first_name": "Natasha",
    "second_name": "Belova", 
    "password": "qwerty"
}
```

Example Response:
``` bash
{
    "id": 3
}
```
#### 1.2 Email verification

Method: `POST`

Endpoint: `auth/verify/email`

Description: Verificates email. After user signed-up, it is required to verificate their email. (Users with non-verified emails cannot sign-in).

Request Parameters:
* `email` (string, requiered)

Response:
* `confirmation code` (string)

Example Request:
``` bash
{
    "email":"dkhfлg@gmail.com"
}
```

Example Response:
``` bash
{
    "confirmationCode": "2266"
}
```
Note: Use confirmation code to verificate user 

#### 1.3. User verification 

Method: `POST`

Endpoint: `auth/verify/user`

Description: Verificates user. When the email is confirmed, you need to notify the server to update the user's status in the database to verified.

Request Parameters:

*`id`(string, required)

Example Request:
``` bash
{
    "id": "3"
}
```
