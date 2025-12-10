-- Add email column to subscriptions table
ALTER TABLE subscriptions ADD COLUMN email VARCHAR(255) NOT NULL DEFAULT '';

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_subscriptions_email ON subscriptions(email);
