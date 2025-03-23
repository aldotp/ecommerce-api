package dto

type SendEmailRequest struct {
	SendEmailType string                 `json:"send_email_type"`      // Type of email, e.g., "activation", "reset_password"
	Email         string                 `json:"email"`                // Recipient's email address
	Subject       string                 `json:"subject"`              // Email subject
	Body          string                 `json:"body"`                 // Email body
	ExtraData     map[string]interface{} `json:"extra_data,omitempty"` // Optional field for additional metadata
}
