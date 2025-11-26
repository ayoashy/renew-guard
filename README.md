# RenewGuard - Subscription Reminder Backend

A production-ready Go backend system that helps users manage their subscriptions and sends automated email notifications when subscriptions are about to expire.

## 🚀 Features

- **User Authentication**: JWT-based authentication with secure password hashing (bcrypt)
- **Subscription Management**: Full CRUD operations for managing subscriptions
- **Smart Notifications**: Automated daily email reminders starting 5 days before expiration
- **Background Scheduler**: Cron-based scheduler for checking and sending notifications
- **Clean Architecture**: Modular structure with repositories, services, and controllers
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **PostgreSQL Database**: Robust data storage with GORM ORM
- **Email Service**: Pluggable email system with SMTP support

## 📋 Table of Contents

- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Development](#development)
- [Deployment](#deployment)

## 🛠 Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Scheduler**: robfig/cron
- **Email**: SMTP
- **Containerization**: Docker & Docker Compose

## 📁 Project Structure

```
renew-guard/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── database/                # Database initialization
│   │   └── database.go
│   ├── models/                  # Data models
│   │   ├── user.go
│   │   ├── subscription.go
│   │   └── notification_log.go
│   ├── repositories/            # Data access layer
│   │   ├── user_repository.go
│   │   ├── subscription_repository.go
│   │   └── notification_log_repository.go
│   ├── services/                # Business logic
│   │   ├── auth_service.go
│   │   ├── subscription_service.go
│   │   └── notification_service.go
│   ├── controllers/             # HTTP handlers
│   │   ├── auth_controller.go
│   │   └── subscription_controller.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   └── error_middleware.go
│   ├── routes/                  # Route definitions
│   │   └── routes.go
│   └── scheduler/               # Background jobs
│       └── scheduler.go
├── pkg/
│   ├── email/                   # Email service
│   │   ├── email.go
│   │   ├── smtp.go
│   │   └── templates.go
│   ├── jwt/                     # JWT utilities
│   │   └── jwt.go
│   └── utils/                   # Helper utilities
│       ├── response.go
│       └── validator.go
├── migrations/                  # SQL migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_subscriptions_table.up.sql
│   ├── 000002_create_subscriptions_table.down.sql
│   ├── 000003_create_notification_logs_table.up.sql
│   └── 000003_create_notification_logs_table.down.sql
├── .env.example                 # Environment variables template
├── Dockerfile                   # Docker image definition
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker & Docker Compose (optional, for containerized setup)
- SMTP credentials for email notifications

### Option 1: Docker Setup (Recommended)

1. **Clone the repository**
   ```bash
   cd renew-guard
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your SMTP credentials
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Check logs**
   ```bash
   docker-compose logs -f app
   ```

The application will be available at `http://localhost:8080`

### Option 2: Local Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   createdb renew_guard
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database and SMTP credentials
   ```

4. **Run the application**
   ```bash
   go run cmd/app/main.go
   ```

   Or use the Makefile:
   ```bash
   make run
   ```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
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

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Subscription Endpoints (Protected)

All subscription endpoints require JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Create Subscription
```http
POST /api/subscriptions
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Netflix",
  "start_date": "2024-01-01T00:00:00Z",
  "duration_days": 30
}
```

#### Get All Subscriptions
```http
GET /api/subscriptions
Authorization: Bearer <token>
```

#### Get Subscription by ID
```http
GET /api/subscriptions/:id
Authorization: Bearer <token>
```

#### Update Subscription
```http
PUT /api/subscriptions/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Netflix Premium",
  "start_date": "2024-01-01T00:00:00Z",
  "duration_days": 30,
  "notification_enabled": true
}
```

#### Delete Subscription
```http
DELETE /api/subscriptions/:id
Authorization: Bearer <token>
```

#### Toggle Notifications
```http
PATCH /api/subscriptions/:id/notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "enabled": false
}
```

### Health Check
```http
GET /health
```

## ⚙️ Configuration

Environment variables are defined in `.env` file:

### Server Configuration
- `SERVER_PORT`: Application port (default: 8080)
- `GIN_MODE`: Gin mode - debug/release (default: debug)
- `APP_ENV`: Application environment (default: development)

### Database Configuration
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: SSL mode (disable/require)

### JWT Configuration
- `JWT_SECRET`: Secret key for JWT signing (change in production!)
- `JWT_EXPIRATION_HOURS`: Token expiration time in hours (default: 72)

### Email Configuration (SMTP)
- `SMTP_HOST`: SMTP server host (e.g., smtp.gmail.com)
- `SMTP_PORT`: SMTP server port (e.g., 587)
- `SMTP_USERNAME`: SMTP username/email
- `SMTP_PASSWORD`: SMTP password/app password
- `SMTP_FROM_EMAIL`: Sender email address
- `SMTP_FROM_NAME`: Sender name

### Scheduler Configuration
- `SCHEDULER_ENABLED`: Enable/disable scheduler (true/false)
- `SCHEDULER_CRON`: Cron expression (default: "0 0 * * *" - midnight daily)
- `NOTIFICATION_DAYS_BEFORE`: Days before expiration to start notifications (default: 5)

## 🔧 Development

### Build the project
```bash
make build
```

### Run tests
```bash
make test
```

### Clean build artifacts
```bash
make clean
```

### Tidy dependencies
```bash
make tidy
```

### Docker commands
```bash
make docker-build    # Build Docker image
make docker-up       # Start containers
make docker-down     # Stop containers
make docker-logs     # View logs
```

## 🌟 How It Works

### Notification System

1. **Scheduler runs daily** at midnight (configurable via `SCHEDULER_CRON`)
2. **Checks all subscriptions** where `notification_enabled = true`
3. **Calculates days until expiration** for each subscription
4. **Sends email notification** if:
   - Days until expiration ≤ 5 (configurable)
   - Days until expiration ≥ 0
   - No notification sent today
5. **Updates `last_notification_sent`** timestamp
6. **Logs notification** in `notification_logs` table

### Email Template

Users receive beautifully formatted HTML emails with:
- Subscription name
- Days remaining
- Expiration date
- Responsive design
- Clear call-to-action

## 🔒 Security Features

- **Password Hashing**: Bcrypt with default cost factor
- **JWT Authentication**: Secure token-based authentication
- **SQL Injection Protection**: GORM parameterized queries
- **CORS Support**: Configurable cross-origin requests
- **Input Validation**: Email and password validation
- **Ownership Verification**: Users can only access their own subscriptions

## 📝 Database Schema

### Users Table
- `id`: Primary key
- `email`: Unique, not null
- `password_hash`: Not null
- `created_at`, `updated_at`: Timestamps

### Subscriptions Table
- `id`: Primary key
- `user_id`: Foreign key to users
- `name`: Subscription name
- `start_date`: Start date
- `duration_days`: Duration in days
- `end_date`: Computed end date
- `notification_enabled`: Boolean flag
- `last_notification_sent`: Last notification timestamp
- `created_at`, `updated_at`: Timestamps

### Notification Logs Table
- `id`: Primary key
- `subscription_id`: Foreign key to subscriptions
- `sent_at`: Timestamp
- `status`: success/failed
- `error_message`: Error details if failed

## 🚢 Deployment

### Production Checklist

1. **Change JWT Secret**: Update `JWT_SECRET` to a strong random value
2. **Configure SMTP**: Set up production email credentials
3. **Set GIN_MODE**: Change to `release`
4. **Enable SSL**: Set `DB_SSLMODE=require` for production database
5. **Secure Environment Variables**: Use secrets management
6. **Set up monitoring**: Add logging and monitoring tools
7. **Configure CORS**: Restrict allowed origins in production

### Docker Production Deploy
```bash
# Build optimized image
docker-compose build

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is open source and available under the MIT License.

## 🆘 Troubleshooting

### Email not sending
- Verify SMTP credentials are correct
- For Gmail, use App Password instead of regular password
- Check firewall allows outbound connections on SMTP port

### Database connection issues
- Verify PostgreSQL is running
- Check database credentials in `.env`
- Ensure database exists

### Build errors
- Run `go mod tidy` to sync dependencies
- Ensure Go version is 1.21 or higher

## 📧 Support

For issues and questions, please create an issue in the repository.

---

Built with ❤️ using Go, Gin, PostgreSQL, and GORM
