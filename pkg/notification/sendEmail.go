package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ecommerce-app/config"
)

type NotificationClient interface {
	SendEmail(email string, subject string, message string) error
}

type notificationClient struct {
	config config.AppConfig
}

func NewNotificationClient(cfg config.AppConfig) NotificationClient {
	return &notificationClient{
		config: cfg,
	}
}

type sender struct {
	Email string `json:"email"`
}

type receiver struct {
	Email string `json:"email"`
}

type emailRequest struct {
	Sender      sender     `json:"sender"`
	To          []receiver `json:"to"`
	Subject     string     `json:"subject"`
	HtmlContent string     `json:"htmlContent"`
}

func (c *notificationClient) SendEmail(
	email string,
	subject string,
	message string,
) error {

	reqBody := emailRequest{
		Sender: sender{
			Email: c.config.EmailFrom,
		},
		To: []receiver{
			{
				Email: email,
			},
		},
		Subject:     subject,
		HtmlContent: fmt.Sprintf("<pre>%s</pre>", message),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", c.config.BrevoAPIKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"failed to send email: %s, response=%s",
			resp.Status,
			string(bodyBytes),
		)
	}

	return nil
}
