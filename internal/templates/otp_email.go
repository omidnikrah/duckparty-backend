package templates

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type OTPEmailData struct {
	OTPCode int
	Year    int
}

const otpEmailHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DuckParty Verification Code</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; padding: 60px 20px;">
        <tr>
            <td align="center">
                <table role="presentation" style="max-width: 500px; width: 100%; border-collapse: collapse; background-color: #ffffff;">
                    <tr>
                        <td style="padding: 60px 40px 40px; text-align: center;">
                            <img src="https://duckparty.s3.eu-north-1.amazonaws.com/ducks/duck-body.png" alt="DuckParty" width="150" height="150" style="width: 150px; height: auto; display: block; margin: 0 auto 30px; border: 0;">
                            <h1 style="margin: 0; color: #f1571f; font-size: 30px; font-weight: bold; letter-spacing: -0.5px;">DuckParty</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 0 40px 40px;">
                            <h2 style="margin: 0 0 12px; color: #000000; font-size: 20px; font-weight: 500; text-align: center;">Your Verification Code</h2>
                            <p style="margin: 0 0 40px; color: #666666; font-size: 15px; line-height: 1.6; text-align: center;">Use this code to verify your account. This code will expire in 2 minutes.</p>
                            
                            <table role="presentation" style="width: 100%; margin: 0 0 40px;">
                                <tr>
                                    <td align="center">
                                        <div style="font-size: 42px; font-weight: 600; color: #000000; letter-spacing: 8px; font-family: 'Courier New', 'Monaco', monospace; padding: 20px 0;">{{.OTPCode}}</div>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 0 0 0; color: #999999; font-size: 13px; line-height: 1.5; text-align: center;">
                                If you didn't request this code, you can safely ignore this email.
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 40px; border-top: 1px solid #e5e5e5; text-align: center;">
                            <p style="margin: 0; color: #999999; font-size: 12px; line-height: 1.5;">
                                © {{.Year}} DuckParty
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const otpEmailTextTemplate = `🦆 DuckParty - Verification Code

Your verification code is: {{.OTPCode}}

Use this code to verify your account. This code will expire in 2 minutes.

If you didn't request this code, you can safely ignore this email.

© {{.Year}} DuckParty`

func GenerateOTPEmailHTML(otpCode int) (string, error) {
	tmpl, err := template.New("otpEmailHTML").Parse(otpEmailHTMLTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	data := OTPEmailData{
		OTPCode: otpCode,
		Year:    time.Now().Year(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

func GenerateOTPEmailText(otpCode int) (string, error) {
	tmpl, err := template.New("otpEmailText").Parse(otpEmailTextTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}

	data := OTPEmailData{
		OTPCode: otpCode,
		Year:    time.Now().Year(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}

	return buf.String(), nil
}
