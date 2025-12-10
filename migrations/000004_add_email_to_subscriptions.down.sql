-- Remove email column from subscriptions table
ALTER TABLE subscriptions DROP COLUMN IF EXISTS email;
DROP INDEX IF EXISTS idx_subscriptions_email;
