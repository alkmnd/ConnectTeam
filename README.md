# ConnectTeam server app 

## Run 
1. First, install go and golang-migrate. For macOS:
``` bash
brew install golang-migrate
```
2. Add .env file in the project.
3. Run docker container. Example:
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

### Содержание

1. [Authentification](#auth)
   
   1.1. [Sign-Up](#sign-up)

   1.2. [Email verification](#verify-email)
   
   1.3. [User verification](#verify-user)
   
   1.4. [Sign-In](#sign-in)

2. [User](#user)

**Note:**

For all REST Api authenticated requests, ensure to include the Authorization header with the value Bearer 
your_access_token. This token is obtained through the user authentication process described in the
first section.
All requests require a valid authorization token in the header.

<a id="auth"></a>
### 1. Authentication
<a id="sign-up"></a>
#### 1.1 Sign-Up

**Method:** `POST`

**Endpoint:** `/auth/sign-up`

**Description:** Registrates user

**Request Parameters:**

* `email` (string, required): User email  

* `phone_number` (string, required): User phone number 

* `first_name` (string, required)

* `second_name` (string, required)

* `password` (string required)

**Response:**

* `id` (string): Signed-up user id

**Example body:**
``` bash
{
    "email": "dkhfлg@gmail.com",
    "phone_number": "89912818155",
    "first_name": "Natasha",
    "second_name": "Belova", 
    "password": "qwerty"
}
```

**Example Response:**
``` bash
{
    "id": 3
}
```
<a id="verify-email"></a>
#### 1.2 Email verification

**Method:** `POST`

**Endpoint:** `auth/verify/email`

**Description:** Verificates email. After user signed-up, it is required to verificate their email. (Users with non-verified emails cannot sign-in).

**Request Parameters:**
* `email` (string, requiered)

**Response:**
* `confirmation code` (string)

**Example Request:**
``` bash
{
    "email":"dkhfлg@gmail.com"
}
```

**Example Response:**
``` bash
{
    "confirmationCode": "2266"
}
```
Note: Use confirmation code to verificate user 

<a id="verify-user"></a>
#### 1.3 User verification

**Method:** `POST`

**Endpoint:** `auth/verify/user`

**Description:** Verificates user. When the email is confirmed, you need to notify the server to update the user's status in the database to verified.

**Request Parameters:**

* `id`(string, required)

**Example Request:**
``` bash
{
    "id": "3"
}
```
<a id="sign-in"></a>
#### 1.4 User Authentication 

**Method:** `POST`

**Endpoint:** `/auth/sign-in/email`

**Description:** Authorizes the user and returns a token for api requests.

**Request Parameters:**

* `email` (string, required)
* `password` (string, required)

**Response:**
* `token` (string)

**Example Request:**
``` bash
{
    "email": "admin@gmail.com",
    "password": "qwert1y"
}
```
<a id="user"></a>
### 2. User 