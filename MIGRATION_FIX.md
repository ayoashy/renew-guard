# Quick Fix for Migration Error

## The Problem

You have existing subscriptions in the database without email values. When GORM tries to add a `NOT NULL` email column, PostgreSQL rejects it because existing rows would have NULL values.

## The Fix

**Changed:** `Email string gorm:"not null"` 
**To:** `Email string gorm:"default:''"`

This allows the migration to succeed by:
1. Adding the email column with empty string default
2. Existing subscriptions get empty string for email
3. New subscriptions will have email populated

## Steps to Apply

### Option 1: Rebuild and Restart (Recommended)
```bash
# Rebuild the Docker image
docker-compose build

# Restart containers
docker-compose down
docker-compose up -d
```

The app should now start successfully!

### Option 2: Alternative - Delete Existing Subscriptions
If you want to start fresh:
```bash
# Stop containers
docker-compose down

# Remove the database volume to start clean
docker volume prune

# Start again
docker-compose up -d
```

## After Fixing

### Existing Subscriptions
- Will have `email: ""` (empty string)
- Notifications won't be sent for these (empty email)
- **Recommended:** Delete and recreate them

### New Subscriptions
- Will have proper email from authenticated user
- Will receive all notifications
- Will get confirmation emails

## Verifying the Fix

```bash
# Check if app is running
docker-compose ps

# Check logs
docker-compose logs -f app

# You should see:
# "Configuration loaded successfully"
# "Database connection established successfully"
# "Auto-migrations completed successfully"
# "Starting server on :8080"
```

## Cleaning Up Old Subscriptions

Once the app is running, you can:

1. **List all subscriptions** (some will have empty email)
2. **Delete old subscriptions** without emails
3. **Create new subscriptions** (will have proper emails)

Or update them programmatically through the database:
```sql
-- Update existing subscriptions with user emails
UPDATE subscriptions s
SET email = (SELECT email FROM users WHERE id = s.user_id)
WHERE email = '';
```

## Summary

✅ Email field now nullable during migration
✅ App will start successfully  
✅ Existing subscriptions will have empty emails (won't get notifications)
✅ New subscriptions will work perfectly
✅ You can manually update old subscriptions or delete them
