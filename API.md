# API Documentation

## Base URL
```
http://localhost:8080
```

## Response Format

All API responses follow this structure:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## Authentication

### Register User

**Endpoint:** `POST /api/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid email format or weak password
- `409 Conflict`: Email already exists
- `500 Internal Server Error`: Server error

---

### Login

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Server error

---

## Subscriptions

All subscription endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Create Subscription

**Endpoint:** `POST /api/subscriptions`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Netflix",
  "start_date": "2024-01-01T00:00:00Z",
  "duration_days": 30
}
```

**Field Descriptions:**
- `name` (required): Name of the subscription (e.g., "Netflix", "Gym Membership")
- `start_date` (required): ISO 8601 formatted start date
- `duration_days` (required): Duration in days (must be > 0)

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Subscription created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30,
    "end_date": "2024-01-31T00:00:00Z",
    "notification_enabled": true,
    "last_notification_sent": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or data
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

### Get All Subscriptions

**Endpoint:** `GET /api/subscriptions`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscriptions retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Netflix",
      "start_date": "2024-01-01T00:00:00Z",
      "duration_days": 30,
      "end_date": "2024-01-31T00:00:00Z",
      "notification_enabled": true,
      "last_notification_sent": "2024-01-26T00:00:00Z",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "name": "Gym Membership",
      "start_date": "2024-01-15T00:00:00Z",
      "duration_days": 365,
      "end_date": "2025-01-15T00:00:00Z",
      "notification_enabled": true,
      "last_notification_sent": null,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

### Get Subscription by ID

**Endpoint:** `GET /api/subscriptions/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**URL Parameters:**
- `id`: Subscription ID

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscription retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30,
    "end_date": "2024-01-31T00:00:00Z",
    "notification_enabled": true,
    "last_notification_sent": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid subscription ID
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Subscription belongs to another user
- `404 Not Found`: Subscription not found
- `500 Internal Server Error`: Server error

---

### Update Subscription

**Endpoint:** `PUT /api/subscriptions/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**
- `id`: Subscription ID

**Request Body:**
```json
{
  "name": "Netflix Premium",
  "start_date": "2024-01-01T00:00:00Z",
  "duration_days": 30,
  "notification_enabled": true
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscription updated successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix Premium",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30,
    "end_date": "2024-01-31T00:00:00Z",
    "notification_enabled": true,
    "last_notification_sent": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or data
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Subscription belongs to another user
- `404 Not Found`: Subscription not found
- `500 Internal Server Error`: Server error

---

### Delete Subscription

**Endpoint:** `DELETE /api/subscriptions/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**URL Parameters:**
- `id`: Subscription ID

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscription deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request`: Invalid subscription ID
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Subscription belongs to another user
- `404 Not Found`: Subscription not found
- `500 Internal Server Error`: Server error

---

### Toggle Notifications

**Endpoint:** `PATCH /api/subscriptions/:id/notifications`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**
- `id`: Subscription ID

**Request Body:**
```json
{
  "enabled": false
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Notification settings updated successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30,
    "end_date": "2024-01-31T00:00:00Z",
    "notification_enabled": false,
    "last_notification_sent": null,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Subscription belongs to another user
- `404 Not Found`: Subscription not found
- `500 Internal Server Error`: Server error

---

## Health Check

### Check API Health

**Endpoint:** `GET /health`

**No authentication required**

**Success Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "renew-guard"
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request succeeded |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Authentication required or invalid token |
| 403 | Forbidden - Authenticated but not authorized |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server error |

---

## Example Usage with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

### Create Subscription
```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30
  }'
```

### Get All Subscriptions
```bash
curl -X GET http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Subscription
```bash
curl -X PUT http://localhost:8080/api/subscriptions/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix Premium",
    "start_date": "2024-01-01T00:00:00Z",
    "duration_days": 30,
    "notification_enabled": true
  }'
```

### Delete Subscription
```bash
curl -X DELETE http://localhost:8080/api/subscriptions/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Toggle Notifications
```bash
curl -X PATCH http://localhost:8080/api/subscriptions/1/notifications \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false
  }'
```
