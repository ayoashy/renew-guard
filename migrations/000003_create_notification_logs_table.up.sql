-- Create notification_logs table
CREATE TABLE IF NOT EXISTS notification_logs (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    error_message TEXT
);

-- Create index for querying logs by subscription and time
CREATE INDEX IF NOT EXISTS idx_notification_logs_subscription_sent ON notification_logs(subscription_id, sent_at);
