A production-ready Go backend system that helps users manage their subscriptions and sends automated email notifications when subscriptions are about to expire.

## Features

- **User Authentication**: JWT-based authentication with secure password hashing (bcrypt)
- **Subscription Management**: Full CRUD operations for managing subscriptions
- **Smart Notifications**: Automated daily email reminders starting 5 days before expiration
- **Background Scheduler**: Cron-based scheduler for checking and sending notifications
- **Clean Architecture**: Modular structure with repositories, services, and controllers
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **PostgreSQL Database**: Robust data storage with GORM ORM
- **Email Service**: Pluggable email system with SMTP support



### Docker Setup

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

### Health Check
```http
GET /health
```
### Docker commands
```bash
make docker-build    # Build Docker image
make docker-up       # Start containers
make docker-down     # Stop containers
make docker-logs     # View logs
```