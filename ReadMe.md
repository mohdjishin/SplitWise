# SplitWise

**_Important:_** **Don't worry about database scripts; GORM will handle the migrations. Start by registering a user.**

## Installation and Setup

### Prerequisites

1. Docker

### Steps

1. Clone the repository
2. Run `docker compose up --build` in the root directory of the project

## Step to run locally without Docker

### Prerequisites

1. Go 
2. PostgreSQL

### Steps

1. Clone the repository
2. Set up PostgreSQL and create a database
3. Update the database configuration in `config.json`
4. incase need to regenerate the swagger file run `swag init` in the root directory of the project
5. Run `go run cmd/main.go` in the root directory of the project

## Usage Examples

## API Documentation
For interactive API documentation, visit [Swagger UI](http://localhost:8080/swagger/index.html) once the application is running.

### Ping 
```bash
curl -I -X GET http://localhost:8080/ping
```

### Register User

```bash
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{
    "email": "foo@example.com",
    "password": "passwOrd@123",
    "name": "foo"
}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{
    "email": "foo@example.com",
    "password": "passwOrd@123"
}'
```

### Create Group with bill

```bash
curl -X POST http://localhost:8080/v1/groups \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
    "groupName": "Weekend Trip",
    "bill": {
        "name": "Trip Expenses",
        "amount": 500.00
    }
}'
```

### Delete group 

```bash
#  not needed as of now. 
curl -X DELETE http://localhost:8080/v1/groups/{groupID} \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN"

```

### List group by owner
```bash
curl -X GET http://localhost:8080/v1/groups \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
```

### Add users to group.

```bash
curl -X POST http://localhost:8080/v1/groups/{groupID}/addMembers \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "userEmailIds": [emailOne, emailTwo, emailthree.....]
}'
```

### Make payment for a group.
```bash
curl -X POST http://localhost:8080/v1/payments \
-H "Content-Type: application/json" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-d '{
    "groupId": <GROUP_ID>,
    "remarks":"remarks"
}'
```
### List Member Groups
```bash

curl -X GET "http://localhost:8080/v1/groups/member-groups?status=PENDING" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-H "Content-Type: application/json"

```

### list pending payments for a user
```bash
curl -X GET http://localhost:8080/v1/payments/pending \
-   H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"
```


### download report based on owner from date to date by default will give last one week data - (today -7)
```bash
curl -X POST http://localhost:8080/v1/report \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "from":"2023-01-01",
    "to":"2024-08-013"
}' \
--output report.pdf
```


Download Report Based on Group ID
```bash
curl -X GET "http://localhost:8080/v1/report/1" \
-H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
--output report.pdf
```


