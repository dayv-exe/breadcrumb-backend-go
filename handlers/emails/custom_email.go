package emails

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func getBody(headerTxt string, promptTxt string, codeParam string) string {
	return fmt.Sprintf(`
<html>
  <body style="margin:0; padding:0; background-color:#f7f7f7; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;">
    <table align="center" width="100%%" cellpadding="0" cellspacing="0" style="max-width:600px; background:white; border-radius:8px; box-shadow:0 2px 8px rgba(0,0,0,0.1); margin-top:40px;">
      <tr>
        <td style="padding:30px 40px;">
          <h2 style="color:#70B8FF; margin:0 0 10px;">🍞 %s</h2>
          <p style="font-size:16px; color:#333;">Hi there,</p>
          <p style="font-size:15px; color:#444;">%s</p>
          <div style="margin:30px 0; text-align:center;">
            <span style="display:inline-block; font-size:24px; font-weight:bold; color:#222; background:#f0f0f0; padding:12px 24px; border-radius:8px;">
              %s
            </span>
          </div>
          <p style="font-size:14px; color:#666;">DO NOT SHARE THIS CODE WITH ANYONE. If you didn't request this, you can safely ignore this email.</p>
          <hr style="margin:30px 0; border:none; border-top:1px solid #eee;">
        </td>
      </tr>
    </table>
  </body>
</html>
`, headerTxt, promptTxt, codeParam)
}

func Handler(ctx context.Context, event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	switch event.TriggerSource {
	case "CustomMessage_SignUp":
		event.Response.EmailSubject = "Verify your Breadcrumb account"
		event.Response.EmailMessage = getBody("Welcome to Breadcrumb!", "Use this code to complete your signup process", event.Request.CodeParameter)

	case "CustomMessage_ForgotPassword":
		event.Response.EmailSubject = "Reset your Breadcrumb password"
		event.Response.EmailMessage = getBody("You requested to reset your password.", "Use this code to complete your password change:", event.Request.CodeParameter)

	case "CustomMessage_UpdateUserAttribute":
		event.Response.EmailSubject = "Confirm your new email address"
		event.Response.EmailMessage = getBody("You're updating your email address.", "Use this code to complete your email address change:", event.Request.CodeParameter)

	default:
		log.Printf("Unknown trigger source: %s", event.TriggerSource)
	}

	return event, nil
}
