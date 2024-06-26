# ConnectTeam server app 

ConnectTeam server app is HTTP-server app accepting HTTP requests from web app (client side of the project) and Game Service (WebSocket server app). Also, ConnectTeam server app sends messages to NotificationService via WebSocket connection.


## Run 
1. First, deploy database (postgres). Use .sql file in ./schema directory to create database.
2. Deploy Redis database.
3. Deploy Notification Service.
4. Add .env and update config.yml files in the project.
5. Run app:
``` bash
go run cmd/main.go
```
### .env 

Example:


``` .env
DB_PASSWORD=qwerty
REDIS_PASSWORD=redis
EMAIL=connectteam66@gmail.com
EMAIL_PASSWORD=pktx fgtm fmwp vjfk
YOOCASSA_API_KEY=test_Mr8jjMiIVtUbPTviWuz8Dv7taCv2SRmnoqxWXey3yvg
SERVICE_API_KEY=YkH8H0LOsoqk6ZZ0MQYkte7proV8Y3QZ
NOTIFICATION_SERVICE_API_KEY=YkH8H0LOsoqk6ZZ0MQYkte7proV8Y3QZ
``` 


* DB_PASSWORD: database password
* REDIS_PASSWORD: redis password
* EMAIL: email for sending messages
* EMAIL_PASSWORD: email password
* YOOCASSA_API_KEY: yoocassa api key
* SERVICE_API_KEY: server api key (used by Game Service)
* NOTIFICATION_SERVICE_API_KEY: Notification Service api key (used by ConnectTeam HTTP-server)



### config.yml
Example:
``` yml
port: "8000"

service_port: "8001"

db:
  username: "postgres"
  host: "localhost"
  port: "5436"
  dbname: "postgres"
  sslmode: "disable"

redis:
  host: "localhost"
  port: "6379"

yookassa:
  shop_id: "371284"

notification_service:
  host: "localhost:8081"
  path: "/ws"

client_origin: "http://localhost:5173"
```
* port: ConnectTeam HTTP-server port (used by client web app)
* service_port: ConnectTeam HTTP-server port (used by Game Service)
* db: database info for connecting
* redis: redis info for connecting
* yookassa.shop_id: shop id from yookassa
* notification_service: notification_service info for connectng
* client_origin: web app domen 


---

## API Documentation

### Contents

1. [Authentification](#auth)
   
   1.1. [Sign-Up](#sign-up)

   1.2. [Email verification](#verify-email)
   
   1.3. [User verification](#verify-user)
   
   1.4. [Sign-In](#sign-in)

   1.5. [Restore Password  (Unauthorized)](#restore-password)

2. [User](#user)
   
   2.1. [Current user](#get-me)
   
   2.2. [Change access](#change-access)

   2.3. [Users list](#users-list)

   2.4. [Password Change](#password-change)

   2.5. [Verify Email On Change](#email-check)

   2.6. [Email Change](#email-change)

   2.7. [Edit Personal Data](#edit-data)

   2.8. [Edit Company Data](#company-change)

   2.9. [Restore Password (Authorized)](#restore-password-auth) 

3. [Plan](#plan)
   
   3.1. [Get My Plan](#get-plan)
   
   3.2. [Select Plan](#new-plan)

   3.3. [Get Users Plans List](#users-plans-lists)

   3.4. [Confirm Plan](#confirm-plan)

   3.5. [Set Plan For User](#set-plan)

   3.6. [Get User Subscriptions](#get-user-subs)

   3.7. [Get Trial](#get-trial)
   
   3.8. [Validate code](#validate-code)

   3.9. [Get Members](#get-members)

   3.10. [Add Member](#add-member)

   3.11. [Delete User From Subscription](#delete-user-from-sub)


4. [Topic](#topic)
   
   4.1. [Create Topic](#create-topic)

   4.2. [Get All Topics](#topics)

   4.3. [Delete Topic](#delete-topic)

   4.4. [Update Topic](#update-topic)

5. [Question](#question)
   
   5.1. [Create Question](#q_create)

   5.2. [Delete Question](q_delete)

   5.3. [Question List](q_list)
   
   5.4. [Update Question](#q_update)

   5.5. [Update Question Tags](#update_tags)

7. [Game](#game)

   6.1. [Create Game](#create-game)

   6.2. [Delete Game](#delete-game)

   6.3. [Get Games Created by User](#get-created)

   6.4. [Get All User Games](#get-games)

   6.5. [Add User to the Game](#add-user-to-game)

   6.6. [Validate Game Invitation Code](#validate-game-code)

   6.7. [Get Game](#gate-game)

8. [WebSocket Server](#ws-server)
   


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
* id (int): Id of the subscription.
* expiry_date (string): Expiration date of the plan.
* holder_id (int): User who owns the plan.
* user_id (int): User who is participant in the plan (for 'basic' and 'advanced' who owns the plan).
* plan_access (string): User access in the plan ('holder' or 'additional')
* plan_type (string): Type of the plan.
* status(string): Plan status ('active', 'on_confirm', 'expired').
* invitation_code (string)


**Example Response:**
``` bash
{
      "id": 1,
    "expiry_date": "2024-02-01T12:34:56Z",
    "holder_id": 1,
    "plan_access": "holder",
    "plan_type": "basic",
    "user_id": 1,
    "status": "on_confirm",
   "invitation_code": "12345d78p0123G56",
   "is_trial": false,
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
* status(string): Status of the plan.

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
      "id": 1,
    "atatus": on_confirm,
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

**Endpoint:** `/plans/active`

**Description:** Returns list of users and their plans.


**Response:**
* data(list)

**Example Response:**
``` bash
{
    "data": [
        {
            "id": 1,
            "plan_type": "premium",
            "user_id": 1,
            "holder_id": 1,
            "expiry_date": "0001-01-01T00:00:00Z",
            "duration": 30,
            "plan_access": "holder",
            "status": "active"
        }
    ]
}
```

**Note:** Allowed for admins only.

<a id="confirm-plan"></a>
#### 3.4. Confirm Plan 

**Method:** `PATCH`

**Endpoint:** `/plans/:id`

**Description:** Sets field 'status' as active, activates user plan.

**URL Parameters:**
* id: Plan id.

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
<a id='get-user-subs'></a>
#### 3.6. Get User Subscriptions

**Method:** `GET`

**Endpoint:** `/plans/`

**Description:** Returns all user subscriptions.

**Response:**

* data(list): list of subscriptions.


**Example Response:**
``` bash
{
    "data": [
        {
            "id": 1,
            "plan_type": "premium",
            "user_id": 1,
            "holder_id": 1,
            "expiry_date": "0001-01-01T00:00:00Z",
            "duration": 30,
            "plan_access": "holder",
            "status": "active"
        }
    ]
}
```

<a id='get-tral'></a>
#### 3.7. Get Trial 

**Method:** `POST`

**Endpoint:** `/plans/trial`

**Description:** Creates user trial plan (subscription).

**Response:**
* plan_type(string): Type of plan the user uses.
* user_id(int): User id.
* holder_id(int): Who owns the plan.
* plan_access(string): User access in the plan ('holder' or 'additional')
* duration(int): The number of days the user can use the plan.
* expiry_date(string):  Expiration date of the plan.
* status(string): Status of the plan.
* is_trial

**Example Response:**
``` bash
{
   "id": 1,
    "atatus": "active",
    "duration": 30,
    "expiry_date": "0001-01-01T00:00:00Z",
    "holder_id": 1,
    "plan_access": "holder",
    "plan_type": "basic",
    "user_id": ,
   "is_trial": true
}
```
<a id='validate-code'></a>
#### 3.8. Validate code

**Method:** `GET`

**Endpoint:** `/validate/plan/:code`

**Description:** Validates code and returns holder id.

**URL Parameters:**
* code(string)

**Response Parameters:**
* id (int): Holder id.


<a id='get-members'></a>
#### 3.9. Get Members

**Method:** `GET`

**Endpoint:** `/plans/members/:code`

**Description:** Returns sunscription members.

**URL Parameters:**
* code(string)

**Response Parameters:**
* list of members.


<a id='get-members'></a>
#### 3.10. Join to Subscription (Plan)

**Method:** `POST`

**Endpoint:** `/plans/join/:code`

**Description:** Adds user to the subscription.

**URL Parameters:**
* code(string)

**Response Parameters:**
* id(int)
* user_id(int)
* holder_id(int)
* plan_type(string)
* plan_access(string)
* status(string)
* duration(int)
* expiry_date(string)


<a id='delete-user-from-sub'></a>
#### 3.11. Delete User From Subscription
**Method:** `DELETE`

**Endpoint:** `/plans/:user_id`

**Description:** Removes user from subscription.

**URL Parameters:**
* user_id(int)

**Response Parameters:**
* status(string): ok if there is no error


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
### 5. Question

<a id='q_create'></a>
#### 5.1. Create Question 

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
#### 5.2. Delete Question

**Method:** `DELETE`

**Endpoint:**  `/questions/:id`

**Description:**  Deletes question. 

**URL Parameters:**
* id(int): Question id.


**Response Parameters:** 
* status (int): ok, if there is no error.


**Example Response:**
```bash
{
   "status": "ok"
}
```

<a id='q_list'></a>
#### 5.3. Question List

**Method:** `GET`

**Endpoint:**  `topics/:topic_id/questions/`

**Description:**  Returns lisy of questions. 

**URL Parameters:**
* topic_id(int): Topic id.


**Response Parameters:** 
* data(list): List of questions.


**Example Response:**
```bash
"data": [
        {
            "id": 1,
            "content": "test1",
            "topic_id": 1,
            "tags": [
                {
                    "Id": "b99c359f-e75e-4c3f-a101-f1d1c6676002",
                    "Name": "Находчивость"
                }
            ]
        },
        {
            "id": 2,
            "content": "test2",
            "topic_id": 1,
            "tags": null
        }
    ]
}
```



<a id='q_update'></a>
#### 5.4. Update Question

**Method:** `PATCH`

**Endpoint:**  `/questions/:id`

**Description:**  Returns lisy of questions. 

**Request Parameters:**
* new_content (string)

**URL Parameters:**
* id(int): Question id.


**Response Parameters:** 
* id (int)
* topic_id (int)
* content (string)

**Example Request Body:**
``` bash
{
    "new_content": "New content???"
}
```

**Example Response:**
``` bash
{
    "id": 2,
    "topic_id": 1,
    "content": "New content???"
}
```

<a id='update_tags'></a>
#### 5.5. Update Question Tags

**Method:** `PUT`

**Endpoint:** `/questions/:id/tags`

**Description:** Updtes question tags.

**Example Request Body:**
```bash
{
    "tags": [
        {"id":"03b08d29-de92-46cc-b944-2be7f7257090"}, 
        {"id":"b99c359f-e75e-4c3f-a101-f1d1c6676002"}, 
        {"id":"03b08d29-de92-46cc-b944-2be7f7257090"}
    ]
}
```

**Example Response Body:**
```bash
{
    "data": [
        {
            "id": "03b08d29-de92-46cc-b944-2be7f7257090",
            "name": "Коммуникабельность"
        }
    ]
}
```


<a id='game'></a>
### 6. Game

<a id='create-game'></a>
#### 6.1. Create Game

**Method:** `POST`

**Endpoint:** `/games/`

**Description:** Creaate new game and returns created game.

**Request Parameters:**
* name(string, required): Game name.
* start_date(string, required): When game starts.

<a id='delete-game'></a>
#### 6.2. Delete Game 

**Method:** `DELETE`

**Endpoint:** `/games/:id`

**Description:** Deletes game by id.

**URL Parameters:**
* id (int): Game id.
 
<a id ='get-created'></a>
#### 6.3. Get Games Created by User

**Method:** `GET`

**Endpoint:** `/games/created/:page`

**Description:** Returns games created by user order by start date.

**Request Parameters:**
* page(int, required)

<a id='get-games'></a>
#### 6.4. Get All User Games

**Method:** `GET`

**Endpoint:** `/games/all/:page`

**Description:** Returns all games that user participated in.

**Request Parameters:**
* page(int, required)

<a id='add-user-to-game'></a>
#### 6.5. Add User to the Game

**Method:**`POST`

**Enpoint:** `/games/:code`

**Description:** Adds user to the game as participant.

**URL Parameters:**
*code(string): Invitation code.


<a id='validate-game-code'></a>
#### 6.6. Validate Game Invitation Code

**Method:** `GET`

**Endpoint:** `validate/game/:code`

**Description:** Checks invitation code and returns game.

**URL Parameters:**
*code(string): Invitation code.

<a id='get-game'></a>
#### 6.7. Get Game

**Method:** `GET`

**Endpoint:** `/games/:id`

**Description:** Returns game date by id.

**URL Parameters:**
* id(int): Game id.

<a id='ws-server'></a>
### 7. Websocket Server 

#### 1. Connect to the Game

**Message to send on server:**
```bash
{
    "action": "join-game", 
    "target": {
        "id":1
    }
}
```

**Broadcasted Message:**
```bash
{
    "action": "join-game",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "not_started",
        "creator_id": 1,
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "Test Test"
            }
        ]
    },
    "sender": {
        "id": 2,
        "name": "Test Test"
    }
}
```

**Received Message:**
```bash
{
    "action": "join-success",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "not_started",
        "creator_id": 1,
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            }
        ]
    },
    "sender": {
        "id": 1,
        "name": "Natasha Belova"
    }
}
```

#### 2. Select Topics for the Game 

**Message to send on server:**
```bash
{
    "action": "select-topic", 
    "target": {
        "id":1
    },
    "sender":{"id":1}, 
    "payload": [ 
            {"id":1}, 
            {"id":2}
        ]
}
```

**Broadcasted Message:**
```bash
{
    "action": "select-topic",
    "payload": [
        {
            "id": 1
        },
        {
            "id": 2
        }
    ],
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "not_started",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": false,
                "title": "Topic1"
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2"
            }
        ],
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": {
        "id": 1,
        "name": "Natasha Belova"
    },
    "time": "0001-01-01T00:00:00Z"
}
```

#### 3. Start Game 

**Message to send on server:**
```bash
{
    "action": "start-game", 
    "target": {
        "id":1
    }
}
```

**Broadcasted Message:**
```bash
{
    "action": "start-game",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "in_progress",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": false,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
```

#### 4. Start Round 

**Message to send on server:**
```bash
{
    "action": "start-round", 
    "target": {
        "id":1
    }, 
    "payload": {
        "id":2
    }
}
```

**Broadcasted Message:**
```bash
{
    "action": "start-round",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "in_progress",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": ""
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": ""
                }
            ],
            "users-questions-left": []
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
```

#### 4. Start Stage (Answer and Rating)

**Message to send on server:**
```bash
{
   "action": "start-stage", 
    "target": {
        "id":1
    }, 
    "sender": {
        "id":1
    }
}
```

**Broadcasted Message (round is not ended):**
```bash
{{
    "action": "start-stage",
    "payload": {
        "number": 1,
        "user": {
            "id": 1,
            "name": "Natasha Belova"
        },
        "question": ""
    },
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "in_progress",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": ""
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": ""
                }
            ],
            "users-questions-left": []
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": {
        "id": 1,
        "name": "Natasha Belova"
    },
    "time": "0001-01-01T00:00:00Z"
}
```

**Broadcasted Message (round is not ended):**
```bash
{{
    "action": "start-stage",
    "payload": {
        "number": 1,
        "user": {
            "id": 1,
            "name": "Natasha Belova"
        },
        "question": ""
    },
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "in_progress",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": ""
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": ""
                }
            ],
            "users-questions-left": []
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": {
        "id": 1,
        "name": "Natasha Belova"
    },
    "time": "0001-01-01T00:00:00Z"
}
```

**Broadcasted Message (round is ended, but game is not):**
```bash
{
    "action": "round-end",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "in_progress",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": false,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": null,
            "users-questions-left": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                }
            ]
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
```

**Broadcasted Message (game is ended, all rounds are ended):**

```bash
{
    "action": "game-end",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "ended",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": null,
            "users-questions-left": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                }
            ]
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
{
    "action": "round-end",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "ended",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": null,
            "users-questions-left": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                }
            ]
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
```

#### 5. Start Answer 

**Message to send on server:**
```bash
{
   "action": "start-answer", 
    "target": {
        "id":1
    }, 
    "sender": {
        "id":1
    }
}
```

**Broadcasted Message:**
```bash
{
    "action": "start-answer",
    "target": {
        "id": 1
    },
    "sender": {
        "id": 2,
        "name": "т е"
    },
    "time": "0001-01-01T00:00:00Z"
}
```

#### 7. End Answer

**Message to send on server:**
```bash
{
   "action": "end-answer", 
    "target": {
        "id":1
    }, 
    "sender": {
        "id":1
    }
}
```

**Broadcasted Message:**
```bash
{
    "action": "end-answer",
    "target": {
        "id": 1
    },
    "sender": {
        "id": 2,
        "name": "т е"
    },
    "time": "0001-01-01T00:00:00Z"
}
```

#### 8. Rate User 

**Message to send on server:**
```bash
{
   "action": "rate-user", 
    "target": {
        "id":1
    }, 
    "sender": {
        "id":1
    }, "payload": {"user_id":2, "value":5}
}
```

**Broadcasted Message:**
```bash
{
    "action": "rate-end",
    "target": {
        "name": "new game",
        "max_size": 3,
        "status": "ended",
        "creator_id": 1,
        "topics": [
            {
                "id": 1,
                "used": true,
                "title": "Topic1",
                "questions": [
                    "",
                    ""
                ]
            },
            {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            }
        ],
        "round": {
            "topic": {
                "id": 2,
                "used": true,
                "title": "Topic2",
                "questions": [
                    "",
                    ""
                ]
            },
            "users-questions": null,
            "users-questions-left": [
                {
                    "number": 1,
                    "user": {
                        "id": 1,
                        "name": "Natasha Belova"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                },
                {
                    "number": 2,
                    "user": {
                        "id": 2,
                        "name": "т е"
                    },
                    "question": "",
                    "rates": [
                        {
                            "value": 5
                        }
                    ]
                }
            ]
        },
        "id": 1,
        "users": [
            {
                "id": 1,
                "name": "Natasha Belova"
            },
            {
                "id": 2,
                "name": "т е"
            }
        ]
    },
    "sender": null,
    "time": "0001-01-01T00:00:00Z"
}
```








