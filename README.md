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

### Contents

1. [Authentification](#auth)
   
   1.1. [Sign-Up](#sign-up)

   1.2. [Email verification](#verify-email)
   
   1.3. [User verification](#verify-user)
   
   1.4. [Sign-In](#sign-in)

2. [User](#user)
   
   2.1. [Current user](#get-me)
   
   2.2. [Change access](#change-access)

   2.3. [Users list](#users-list)

   2.4. [Password Change](#password-change)

   2.5. [Verify Email On Change](#email-check)

   2.6. [Change Email](#email-change)

   2.7. [Edit Personal Data](#edit-data)

   


**Note:**

For all REST Api authenticated requests, ensure to include the Authorization header with the value Bearer 
your_access_token. This token is obtained through the user authentication process described in the
first section.
All requests require a valid authorization token in the header.

<a id="auth"></a>
### 1. Authentication
<a id="sign-up"></a>
#### 1.1. Sign-Up

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

**Example Request Body:**
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
#### 1.2. Email verification

**Method:** `POST`

**Endpoint:** `auth/verify-email`

**Description:** Verificates email. After user signed-up, it is required to verificate their email. (Users with non-verified emails cannot sign-in).

**Request Parameters:**
* `email` (string, requiered)

**Response:**

* `id` (int)

**Example Request Body:**
``` bash
{
    "email":"dkhfлg@gmail.com"
}
```

**Example Response:**
``` bash
{
    "id": 1
}
```
Note: Use confirmation code to verificate user 

<a id="verify-user"></a>
#### 1.3. User verification

**Method:** `POST`

**Endpoint:** `auth/verify-user`

**Description:** Verificates user. When the email is confirmed, you need to notify the server to update the user's status in the database to verified.

**Request Parameters:**

* `id`(string, required): User id
* `code`(string, required): User-entered verification code.

**Example Request Body:**
``` bash
{
    "id":  "1", 
    "code": "8956"
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

**Example Request body:**
``` bash
{
    "email": "admin@gmail.com",
    "password": "qwert1y"
}
```

**Example Response:**
```bash
{  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY1MTcyNzksImlhdCI6MTcwNjQ3NDA3OSwidXNlcl9pZCI6MSwiYWNjZXNzIoidXNlciJ9.SlEf62NQ84OuXQ6royCsjfZwzzc7iOmJs5HlgtoXAuY"
}
```

<a id="user"></a>
### 2. User 

<a id="get-me"></a>
#### 2.1. Current user 

**Method:** `GET`

**Endpoint:** `/users/me`

**Description:** `Returns information about the user.`

**Example Response:**
```bash
{
    "access": "user",
    "comppany_name": "",
    "email": "admin@gmail.com",
    "first_name": "Natasha",
    "id": 1,
    "image": "",
    "second_name": "Belova"
}
```

<a id="change-access"></a>
#### 2.2. Change access

**Method:** `PATCH`

**Endpoint:** `/users/change-access`

**Description:** Changes authenticated user access if current user is admin.

**Request Parameters:** 
* id (string, required): User id for changing access.
* access (string, required): One of the user roles in the system that needs to be changed to.

**Example Request body:**
```bash
{
    "id": "1", 
    "access": "admin"
}
```
<a id="users-list"></a>
#### 2.3. Users list 

**Method:** `GET`

**Endpoint:** `users/list`

**Description:** Returns list of users.

**Example Response:**
```bash
{
    "data": [
        {
            "id": 2,
            "email": "admin@gmail.com",
            "first_name": "q",
            "second_name": "q",
            "access": "admin",
            "company_name": "",
            "profile_image": ""
        },
        {
            "id": 1,
            "email": "b@gmail.com",
            "first_name": "fg",
            "second_name": "cdvf",
            "access": "admin",
            "company_name": "",
            "profile_image": ""
        }
    ]
}
```
<a id="password-change"></a>
#### 2.4. Password Change 

**Method:** `PATCH`

**Enpoint:** `/users/change-password`

**Description:** Changes user password.

**Request Parameters:** 
* new_password (string, required): New password
* old_password (string, required): Old password

**Example Request Body:** 
```bash
{
    "new_password": "qwerty",
    "old_password": "qwerty1"
}
```
<a id="email-check"></a>
#### 2.5. Verify Email On Change
**Method:** `POST`

**Endpoint:** `users/verify-email`

**Description:** Checks email when and password user changes email and send verification code.

**Request Parameters:**
* email (string, required): User new email.

**Example Request Body:**
```bash
{
    "email":"ivandoronin22@gmail.com", 
    "password":"qwerty"
}
```

<a id="email-change"></a>
#### 2.6. Email Change

**Method:** `PUTCH`

**Endpoint:** `users/change-email`

**Description:** Changes user email if verification code is valid.

**Request Parameters:**
* new_email (string, required): User new email.
* code (string, required): Verification code sent by user

**Example Request Body:**
```bash
{
    "new_email":"ivandoronin22@gmail.com",
    "code": "6180"
}
```
<a id="edit-data"></a>
#### 2.7. Edit Personal Data

**Method:** `PUTCH`

**Endpoint:** `/users/edit-info`

**Description:** Changes user's first name, second name, description

**Request Parameters:**

* first_name(string, required)
* second_name(string, required)
* description(string, required)


**Example Request Body:**
``` bash
{
    "first_name":"Natasha",
    "second_name": "Belova",
    "description": ""
}

```


<a id="company-change"></a>
#### 2.5. Company Change 








**Method:** ``

**Endpoint:** `/`

**Description:** 

**Request Parameters:**

* 

**Response:**

* 

**Example Body:**
``` bash

```

**Example Response:**
``` bash

```
