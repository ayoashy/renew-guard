package email

import (
	"fmt"
	"time"
)

// GetSubscriptionConfirmationTemplate generates HTML email template for subscription confirmation
func GetSubscriptionConfirmationTemplate(subscriptionName string, startDate time.Time, endDate time.Time) string {
	durationDays := int(endDate.Sub(startDate).Hours() / 24)
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 40px 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .content {
            padding: 30px;
        }
        .success-icon {
            text-align: center;
            font-size: 60px;
            margin: 20px 0;
        }
        .info-box {
            background: #f0f7ff;
            border-left: 4px solid #667eea;
            padding: 20px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .info-box h3 {
            margin-top: 0;
            color: #667eea;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            padding: 10px 0;
            border-bottom: 1px solid #eee;
        }
        .detail-label {
            font-weight: bold;
            color: #666;
        }
        .detail-value {
            color: #333;
        }
        .highlight {
            background: #fff3cd;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            text-align: center;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #666;
            font-size: 12px;
            background: #f9f9f9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Subscription Created!</h1>
        </div>
        <div class="content">
            <div class="success-icon">🎉</div>
            
            <p>Great news! Your subscription has been successfully added to RenewGuard.</p>
            
            <div class="info-box">
                <h3>📋 Subscription Details</h3>
                <div class="detail-row">
                    <span class="detail-label">Service Name:</span>
                    <span class="detail-value">%s</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Start Date:</span>
                    <span class="detail-value">%s</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Duration:</span>
                    <span class="detail-value">%d days</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Expires On:</span>
                    <span class="detail-value">%s</span>
                </div>
            </div>
            
            <div class="highlight">
                <strong>🔔 Notifications Enabled</strong><br>
                We'll remind you 5 days before expiration
            </div>
            
            <p><strong>What happens next?</strong></p>
            <ul>
                <li>✅ Your subscription is now being tracked</li>
                <li>📧 You'll receive daily reminders starting 5 days before expiration</li>
                <li>⚙️ You can manage notification settings anytime</li>
                <li>📊 Monitor all your subscriptions in one place</li>
            </ul>
            
            <p style="margin-top: 30px;">
                <strong>Need to make changes?</strong><br>
                You can update or delete this subscription anytime through the RenewGuard dashboard.
            </p>
        </div>
        <div class="footer">
            <p>This is a confirmation email from RenewGuard</p>
            <p>Subscription Reminder System</p>
            <p>Never miss a renewal date again! 🎯</p>
        </div>
    </div>
</body>
</html>
`, subscriptionName, startDate.Format("Monday, January 2, 2006"), durationDays, endDate.Format("Monday, January 2, 2006"))
}

// GetSubscriptionConfirmationSubject generates subject for confirmation email
func GetSubscriptionConfirmationSubject(subscriptionName string) string {
	return fmt.Sprintf("✅ %s subscription added to RenewGuard", subscriptionName)
}
