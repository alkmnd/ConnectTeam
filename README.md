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

   1.5. [Restore Password  (Unauthorized)](#restore-password)

3. [User](#user)
   
   2.1. [Current user](#get-me)
   
   2.2. [Change access](#change-access)

   2.3. [Users list](#users-list)

   2.4. [Password Change](#password-change)

   2.5. [Verify Email On Change](#email-check)

   2.6. [Email Change](#email-change)

   2.7. [Edit Personal Data](#edit-data)

   2.8. [Edit Company Data](#company-change)

   2.9. [Restore Password (Authorized)](#restore-password-auth) 

4. [Plan](#plan)
   
   3.1. [Get My Plan](#get-plan)
   
   3.2. [Select Plan](#new-plan)

   3.3. [Get Users Plans List](#users-plans-lists)

   3.4. [Confirm Plan](#confirm-plan)

   3.5. [Set Plan For User](#set-plan)

5. [Topic](#topic)
   
   5.1. [Create Topic](#create-topic)

   5.2. [Get All Topics](#topics)

   5.3. [Delete Topic](#delete-topic)

   5.4. [Update Topic](#update-topic)

6. [Question](#question)
   6.1. [Create Question](#q_create)
   
   6.2. [Update Question](#q_update)

   6.3. [Delete Question](q_delete)

   6.4. [Question List](q_list)
   


**Note:**

For all REST Api authenticated requests, ensure to include the Authorization header with the value Bearer 
<your_access_token>. This token is obtained through the user authentication process described in the
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
* `access` (string): User access.

**Example Request body:**
``` bash
{
    "email": "admin@gmail.com",
    "password": "qwert1y"
}
```

**Example Response:**
```bash
{  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY1MTcyNzksImlhdCI6MTcwNjQ3NDA3OSwidXNlcl9pZCI6MSwiYWNjZXNzIoidXNlciJ9.SlEf62NQ84OuXQ6royCsjfZwzzc7iOmJs5HlgtoXAuY",
"access": "user"
}
```
<a id="restore-password"></a>
#### 1.5. Restore Password (Unauthorized)

**Method:** `PATCH`

**Endpoint:** `/auth/password`

**Description:** Sends new password on email

**Request Parameters:**
* email(string, required)

**Response Parameters:**
* status(string): ok, if there is no error

**Example Request Body:** 
```bash
{
    "email":"blvantla@gmail.com"
}
```

**Example Response:**
```bash
{
   "status": "ok"
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
    "company_info": "o_o",
    "company_logo": "",
    "company_url": "0_o",
    "company_name": "Yandex",
    "description": "=)",
    "email": "blvantla@gmail.com",
    "first_name": "Ksusha",
    "id": 1,
    "profile_image": "",
    "second_name": "Belova"
}
```

<a id="change-access"></a>
#### 2.2. Change access

**Method:** `PATCH`

**Endpoint:** `/users/access`

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
            "profile_image": "",
            "company_info": "o_o",
            "company_logo": "",
            "company_url": "0_o",
            "description": "=)",
        },
        {
            "id": 1,
            "email": "b@gmail.com",
            "first_name": "fg",
            "second_name": "cdvf",
            "access": "admin",
            "company_name": "",
            "profile_image": "",
            "company_info": "o_o",
             "company_logo": "",
             "company_url": "0_o",
             "description": "=)",
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

**Method:** `PATCH`

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

**Method:** `PATCH`

**Endpoint:** `/users/info`

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
#### 2.8. Edit Company Data 

**Method:** `PATCH`

**Endpoint:** `/users/company`

**Description:** Changes user's company data (company name, company info, company web-site)

**Request Parameters:**

* company_name(string, required)
* company_info(string, required)
* company_url(string, required)


**Example Request Body:**
``` bash
{
    "company_name":"Yandex",
    "company_info": "o_o",
    "company_url": "0_o"
}

```
<a id="restore-password-auth"> </a>
#### 2.9. Restore Password (Authorized)

**Method:** `GET`

**Endpoint:** `/users/password`

**Description:** Resets the user password and sends message on the user email with new password.

**Response:**

* status(string): "ok", if there is no error.

**Example Response:**
```bash
{
    "status": "ok"
}
```

<a id="plan"></a>
### 3. Plan

<a id="get-plan"></a>
#### 3.1. Get My Plan 

**Method:** `GET`

**Endpoint:** `/plans/current`

**Description:** Returns current user's plan

**Response:**

* expiry_date (string): Expiration date of the plan.
* holder_id (int): User who owns the plan.
* user_id (int): User who is participant in the plan (for 'basic' and 'advanced' who owns the plan).
* plan_access (string): User access in the plan ('holder' or 'additional')
* plan_type (string): Type of the plan.


**Example Response:**
``` bash
{
    "expiry_date": "2024-02-01T12:34:56Z",
    "holder_id": 1,
    "plan_access": "holder",
    "plan_type": "basic",
    "user_id": 1,
    "confirmed": true
}
```
<a id="new-plan"></a>
#### 3.2. Select Plan 

**Method:** `POST`

**Endpoint:** `/plans/purchase`

**Description:** Adds (or updates) a record about the plan for the user currently in the system (using token).

**Request Parameters:**

* duration(int, required):  How long  the user wants to use the plan.
* plan_type(string, required): Type of plan the user wants to use.

**Response:**

* plan_type(string): Type of plan the user uses.
* user_id(int): User id.
* holder_id(int): Who owns the plan.
* plan_access(string): User access in the plan ('holder' or 'additional')
* duration(int): The number of days the user can use the plan.
* expiry_date(string):  Expiration date of the plan.
* confirmed(bool): If admin confirmed the purchase.

**Example Body:**
``` bash
{
    "duration": 30, 
    "plan_type": "premium"
}
```

**Example Response:**
``` bash
{
    "confirmed": false,
    "duration": 30,
    "expiry_date": "0001-01-01T00:00:00Z",
    "holder_id": 1,
    "plan_access": "holder",
    "plan_type": "basic",
    "user_id": 1
}
```


<a id="users-plans-lists"></a>
#### 3.3. Get Users Plans List 

**Method:** `GET`

**Endpoint:** `/plans/`

**Description:** Returns list of users and their plans.


**Response:**
* data(list)

**Example Response:**
``` bash
{
    "data": [
        {
            "plan_type": "premium",
            "user_id": 1,
            "holder_id": 1,
            "expiry_date": "0001-01-01T00:00:00Z",
            "duration": 30,
            "plan_access": "holder",
            "confirmed": false
        }
    ]
}
```

**Note:** Allowed for admins only.

<a id="confirm-plan"></a>
#### 3.4. Confirm Plan 

**Method:** `PATCH`

**Endpoint:** `/plans/:id`

**Description:** Sets field 'confirmed' true, activates user plan.

**URL Parameters:**
* id: Plan id (equal to user id)

**Response:**

* status (string): Returns status ok if there is no error.


**Example Response:**
``` bash
{
    "status": "ok"
}

```

**Note:** Allowed for admins only.

<a id='set-plan'></a>
#### 3.5. Set Plan For User 

**Note**: Allowed for admins only

**Method:** `POST`

**Endpoint:** `/plans/:user_id`

**Description:** Sets new plan to the user by id.

**URL Parameters:**

* user_id: User Id to which the plan is assigned.

**Request Parameters:**

* expiry_date(string, required): Plan expiration date according to RFC3339 standard.
* plan_type(string, required): Type of plan the user is assigned.

**Response:**

* status (string): Returns status ok if there is no error.


**Example Response:**
``` bash
{
    "status": "ok"
}
```
<a id='topic'></a>
### 5. Topic

<a id='create-topic'></a>
#### 5.1. Create Topic

**Note**: Allowed for admins only

**Method:** `POST`

**Endpoint:** `/topics/`

**Description:** Creates new topic

**Request Parameters:**
* title(string, required): The title of a new topic.

**Response:**
* id(int): Id of the newly created topic.

**Example Request Body:**
``` bash
{
    "title": "New Topic"
}
```

**Example Response:**
``` bash
{
   "id": 1
}
```

<a id='topics'></a>
#### 5.2. Get All Topics

**Method:** `GET`

**Endpoint:** `/topics/`

**Description:** Returns list of all topics 

**Response:**
* data([]): List of items.

**Example Response:**
``` bash
{
    "data": [
        {
            "id": 1,
            "title": "h"
        }
    ]
}
```

<a id='delete-topic'></a>
#### 5.3. Delete Topic

**Method:** `DELETE`

**Endpoint:** `/topics/:id`

**Description:** Deletes topic by id.

**URL Parameters:**
* id(int): Id of the topic.

**Response:**
status(string): "ok" if there is no error.

**Example Response:**

``` bash
{
    "status": "ok"
}
```
<a id='update-topic'></a>
#### 5.4. Update Topic 

**Method:** `PATCH`

**Endpoint:** `/plans/:id`

**Description:** Update the topic title.

**URL Parameters:**
* id(int): Id of the uodated topic.

**Request Parameters:**
* title(string, required): new title of the topic.

**Response:**
* status(string): "ok", if there is no error.

**Example Request Body:**
```bash
{
   "title": "New Title"
}
```

**Example Response:**
``` bash
{
    "status": "ok"
}
```
<a id='question'></a>
### 6. Question

<a id='q_create'></a>
#### 6.1. Create Question 

**Method:** `POST`

**Endpoint:**  `topics/:topic_id/questions/`

**Description:**  Creates new question in the topic with topic_id. 

**URL Parameters:**
* topic_id(int): Topic id.

**Request Parameters**

* content (string): Content of the question.

**Response Parameters:** 
* id (int): Id of the created question.

**Example Request Body:**
```bash
{
   "content":"New question?"
}
```

**Example Response:**
```bash
{
   "id": 1
}
```

<a id='q_delete'></a>
#### 6.2. Delete Question

**Method:** `DELETE`

**Endpoint:**  `topics/:topic_id/questions/:id`

**Description:**  Deletes question. 

**URL Parameters:**
* topic_id(int): Topic id.
* id(int): Question id.


**Response Parameters:** 
* status (int): ok, if there is no error.


**Example Response:**
```bash
{
   "status": "ok"
}
```
