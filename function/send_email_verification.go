package function

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailVerification(digitCode string, email string) string {
	// Templates for email code verification
	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_ACC"),
		os.Getenv("SMTP_ACC_PASSWORD"),
		"smtp.gmail.com",
	)
	msg := "Subject: Your Verification Code\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html>" +
		"<body style=\"font-family: Arial, sans-serif; margin: 0; padding: 0;\">" +
		"<div style=\"background-color: #f2f2f2; padding: 20px;\">" +
		"<div style=\"background-color: #ffffff; max-width: 600px; margin: 0 auto; padding: 20px; border-radius: 10px;\">" +
		"<h2 style=\"text-align: center; color: #4CAF50;\">Verification Code</h2>" +
		"<p style=\"text-align: center; font-size: 16px;\">Please use the following code to verify your email address:</p>" +
		"<h3 style=\"text-align: center; font-size: 24px; color: #333;\">" + digitCode + "</h3>" +
		"<p style=\"text-align: center; color: #777;\">If you did not request this code, please ignore this email.</p>" +
		"<hr style=\"border: none; border-top: 1px solid #ddd; margin: 20px 0;\" />" +
		"<p style=\"text-align: center; font-size: 12px; color: #aaa;\">&copy; 2024 Betamart. All rights reserved.</p>" +
		"</div>" +
		"</div>" +
		"</body>" +
		"</html>"

	// Send email to user
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"betamartauth@gmail.com",
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
		return fmt.Sprintln(err)
	}

	return "Success"
}
