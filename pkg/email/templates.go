package email

import (
	"fmt"
	"time"
)

// GetExpirationWarningTemplate generates HTML email template for subscription expiration warning
func GetExpirationWarningTemplate(subscriptionName string, daysLeft int, endDate time.Time) string {
	htmlBody := fmt.Sprintf(`
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
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px;
            border-radius: 10px 10px 0 0;
            text-align: center;
        }
        .content {
            background: #f9f9f9;
            padding: 30px;
            border-radius: 0 0 10px 10px;
        }
        .warning-box {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .info-box {
            background: white;
            padding: 20px;
            margin: 20px 0;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #666;
            font-size: 12px;
        }
        .highlight {
            color: #667eea;
            font-weight: bold;
            font-size: 24px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔔 Subscription Expiration Warning</h1>
    </div>
    <div class="content">
        <div class="warning-box">
            <h2 style="margin-top: 0;">⏰ Action Required</h2>
            <p>Your subscription is expiring soon!</p>
        </div>
        
        <div class="info-box">
            <h3>Subscription Details:</h3>
            <p><strong>Service:</strong> %s</p>
            <p><strong>Days Remaining:</strong> <span class="highlight">%d</span></p>
            <p><strong>Expiration Date:</strong> %s</p>
        </div>
        
        <p>Don't forget to renew your subscription to continue enjoying uninterrupted service.</p>
        
        <p>If you've already renewed, you can safely ignore this message.</p>
    </div>
    <div class="footer">
        <p>This is an automated notification from RenewGuard</p>
        <p>You're receiving this because you enabled notifications for this subscription</p>
    </div>
</body>
</html>
`, subscriptionName, daysLeft, endDate.Format("Monday, January 2, 2006"))

	return htmlBody
}

// GetExpirationWarningSubject generates email subject for expiration warning
func GetExpirationWarningSubject(subscriptionName string, daysLeft int) string {
	if daysLeft == 0 {
		return fmt.Sprintf("🚨 URGENT: Your %s subscription expires TODAY!", subscriptionName)
	} else if daysLeft == 1 {
		return fmt.Sprintf("⚠️ Your %s subscription expires TOMORROW!", subscriptionName)
	}
	return fmt.Sprintf("⚠️ Your %s subscription expires in %d days", subscriptionName, daysLeft)
}
