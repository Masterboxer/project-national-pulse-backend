# 📘 User API - Usage Guide

This API provides basic CRUD operations for managing users in the system.

## 📍 Base URL

```
http://localhost:8000
```

Make sure your server is running on this address. Adjust the port if needed.

---

## 🧑‍💻 Endpoints

### 1. 📥 Create a New User

**POST** `/users`

Creates a new user with a name and email.

#### Request:

```bash
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sriharish",
    "email": "sriharish@google.com"
    "password": "Test@123"
}'
```

#### Response:

```json
{
  "id": 1,
  "name": "Sriharish",
  "email": "sriharish@google.com"
}
```

---

### 2. 📤 Get All Users

**GET** `/users`

Fetches all users.

#### Request:

```bash
curl http://localhost:8000/users
```

#### Response:

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "johndoe@example.com"
  },
  ...
]
```

---

### 3. 🔍 Get a User by ID

**GET** `/users/{id}`

Fetches a specific user by ID.

#### Request:

```bash
curl http://localhost:8000/users/1
```

#### Response:

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```

If the user is not found:

```http
HTTP/1.1 404 Not Found
```

---

### 4. ✏️ Update a User

**PUT** `/users/{id}`

Updates a user's name and email.

#### Request:

```bash
curl -X PUT http://localhost:8000/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "johnupdated@example.com"
}'
```

#### Response:

```json
{
  "id": 1,
  "name": "John Updated",
  "email": "johnupdated@example.com"
}
```

---

### 5. ❌ Delete a User

**DELETE** `/users/{id}`

Deletes a user by ID.

#### Request:

```bash
curl -X DELETE http://localhost:8000/users/1
```

#### Response:

```json
"User deleted"
```

If the user is not found:

```http
HTTP/1.1 404 Not Found
```

---

## 🛠️ Requirements

- Go server must be running with the endpoints correctly wired.
- PostgreSQL database must be up and reachable.
- Users table should exist with columns: `id`, `name`, `email`.

---

## 💪 Run with Docker

To run the application and PostgreSQL database using Docker:

### 1. Ensure the following files are in your project root:

- `Dockerfile`
- `docker-compose.yml`

### 2. Build and start the containers:

```bash
sudo docker compose up --build
```

This will:

- Build the Go app container (`civic-action-app`)
- Start a PostgreSQL container (`civic_action_db`) with database `pulse`
- Expose the API on port `8000`

### 2.5. If Postgres file path is needed to be added temporarily

```
$env:Path += ";C:\Program Files\PostgreSQL\15\bin"
```

### 3. API Access

Once running, the API will be available at:

```
http://localhost:8000
```

You can now use the cURL commands above to test the endpoints.

---

## 📌 Notes

- All responses are in JSON format.
- No authentication is currently implemented.
- Make sure to use unique emails when creating users if your database has a unique constraint on the `email` field.

---
