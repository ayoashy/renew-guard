package controllers

import (
	"net/http"
	"renew-guard/pkg/email"
	"renew-guard/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EmailTestController struct {
	emailService email.EmailService
}

func NewEmailTestController(emailService email.EmailService) *EmailTestController {
	return &EmailTestController{
		emailService: emailService,
	}
}

type TestEmailRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// SendTestEmail sends a test email to verify SMTP configuration
// @Summary Send test email
// @Tags email-test
// @Accept json
// @Produce json
// @Param request body TestEmailRequest true "Test email details"
// @Success 200 {object} Response
// @Router /api/test/email [post]
func (ctrl *EmailTestController) SendTestEmail(c *gin.Context) {
	var req TestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body. Email and name are required.")
		return
	}

	// Validate email format
	if !utils.IsValidEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Create test email content
	subject := "🧪 RenewGuard SMTP Test Email"
	htmlBody := getTestEmailTemplate(req.Name)

	// Send email asynchronously (don't block the response)
	go func() {
		err := ctrl.emailService.SendHTML(req.Email, subject, htmlBody)
		if err != nil {
			// Log error but don't fail the request since it's already responded
			// In production, you might want to log this to a monitoring system
			println("Failed to send test email:", err.Error())
		} else {
			println("Test email sent successfully to:", req.Email)
		}
	}()

	// Respond immediately
	utils.SuccessResponse(c, http.StatusOK, "Test email queued for delivery!", gin.H{
		"recipient": req.Email,
		"name":      req.Name,
		"note":      "Email is being sent in the background. Check your inbox in a few moments.",
	})
}

// getTestEmailTemplate creates a test email template
func getTestEmailTemplate(name string) string {
	return `
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
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .info-box h3 {
            margin-top: 0;
            color: #667eea;
        }
        .footer {
            text-align: center;
            padding: 20px;
            color: #666;
            font-size: 12px;
            background: #f9f9f9;
        }
        .button {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 SMTP Test Successful!</h1>
        </div>
        <div class="content">
            <div class="success-icon">✅</div>
            
            <h2>Hello, ` + name + `!</h2>
            
            <p>Great news! Your SMTP configuration is working perfectly.</p>
            
            <div class="info-box">
                <h3>✨ What This Means</h3>
                <ul>
                    <li>✅ SMTP server connection successful</li>
                    <li>✅ Authentication working correctly</li>
                    <li>✅ Email delivery functioning</li>
                    <li>✅ HTML formatting supported</li>
                </ul>
            </div>
            
            <p><strong>Your email configuration details:</strong></p>
            <ul>
                <li><strong>Recipient:</strong> You (` + name + `)</li>
                <li><strong>Email Status:</strong> Successfully delivered</li>
                <li><strong>Email Type:</strong> HTML formatted</li>
            </ul>
            
            <p>Your RenewGuard subscription reminder system is now ready to send notification emails!</p>
            
            <div style="text-align: center;">
                <p style="color: #667eea; font-size: 18px; font-weight: bold;">
                    🎉 You're all set!
                </p>
            </div>
        </div>
        <div class="footer">
            <p>This is a test email from RenewGuard</p>
            <p>Subscription Reminder Backend System</p>
            <p>Built with Go, Gin, PostgreSQL & ❤️</p>
        </div>
    </div>
</body>
</html>
`
}
