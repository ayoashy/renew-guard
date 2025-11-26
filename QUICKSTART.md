# 🚀 RenewGuard - Quick Start Guide

Get started with the subscription reminder backend in under 5 minutes!

## Prerequisites

- Docker & Docker Compose (easiest)
- OR Go 1.21+ and PostgreSQL 15+

## Quick Start with Docker (Recommended)

### 1. Configure Email Settings

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and update these SMTP settings:
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password-here
```

**For Gmail:** Generate an App Password at https://myaccount.google.com/apppasswords

### 2. Start the Application

```bash
docker-compose up -d
```

That's it! The application is now running at http://localhost:8080

### 3. Check Status

```bash
# View logs
docker-compose logs -f app

# Check if containers are running
docker-compose ps
```

## Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `token` from the response!

### 2. Create a Subscription

```bash
# Replace YOUR_TOKEN with the token from registration
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix",
    "start_date": "2024-12-01T00:00:00Z",
    "duration_days": 7
  }'
```

### 3. List Your Subscriptions

```bash
curl -X GET http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## What Happens Next?

- The **scheduler runs daily at midnight** (configurable)
- It checks for subscriptions expiring within **5 days**
- Sends **email notifications** daily until expiration
- Logs all notification attempts in the database

## Useful Commands

```bash
# Stop the application
docker-compose down

# View app logs
docker-compose logs -f app

# View database logs
docker-compose logs -f postgres

# Restart the app
docker-compose restart app

# Rebuild after code changes
docker-compose build && docker-compose up -d
```

## Local Development (Without Docker)

### 1. Setup Database

```bash
# Create PostgreSQL database
createdb renew_guard
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your local database and SMTP settings
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Application

```bash
go run cmd/app/main.go
```

Or use the Makefile:

```bash
make run
```

## Environment Variables Reference

### Required Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=renew_guard

# JWT Secret (change this!)
JWT_SECRET=your-super-secret-key

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### Optional Variables

```env
SERVER_PORT=8080              # Default: 8080
GIN_MODE=release              # debug or release
JWT_EXPIRATION_HOURS=72       # Token validity in hours
SCHEDULER_ENABLED=true        # Enable background scheduler
SCHEDULER_CRON="0 0 * * *"   # Run at midnight daily
NOTIFICATION_DAYS_BEFORE=5    # Warning days before expiration
```

## API Endpoints

### Authentication (Public)
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token

### Subscriptions (Protected - Requires JWT)
- `POST /api/subscriptions` - Create subscription
- `GET /api/subscriptions` - List all subscriptions
- `GET /api/subscriptions/:id` - Get subscription by ID
- `PUT /api/subscriptions/:id` - Update subscription
- `DELETE /api/subscriptions/:id` - Delete subscription
- `PATCH /api/subscriptions/:id/notifications` - Toggle notifications

### Health Check
- `GET /health` - Check API status

## Troubleshooting

### Emails not sending?
1. Check SMTP credentials in `.env`
2. For Gmail, use App Password (not regular password)
3. Check firewall allows port 587
4. View logs: `docker-compose logs -f app`

### Database connection error?
1. Ensure PostgreSQL is running: `docker-compose ps`
2. Check credentials in `.env`
3. Verify database exists

### Build errors?
```bash
# Clean and rebuild
make clean
make build

# Or with Docker
docker-compose build --no-cache
```

## Next Steps

1. ✅ Test all API endpoints
2. ✅ Create test subscriptions with different expiration dates
3. ✅ Monitor scheduler logs for notification runs
4. ✅ Check email inbox for notification emails
5. ✅ Review the full documentation in `README.md` and `API.md`

## Production Deployment

Before deploying to production:

1. **Change JWT Secret**: Generate a strong random key
2. **Set GIN_MODE**: Change to `release`
3. **Enable DB SSL**: Set `DB_SSLMODE=require`
4. **Use Secrets Management**: Don't commit `.env` to git
5. **Set up Monitoring**: Add logging and alerting
6. **Configure CORS**: Restrict allowed origins

## Support

- 📖 Full Documentation: See `README.md`
- 🔌 API Reference: See `API.md`
- 🐛 Issues: Report in project repository

---

**Happy Coding! 🎉**

Built with ❤️ using Go, Gin, PostgreSQL, and Docker
