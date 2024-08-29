package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// sendWhatsAppNotification sends a WhatsApp message via Twilio
func sendWhatsAppNotification(toPhoneNumber string, messageBody string) error {
	twilioAccountSID := os.Getenv("twilioAccountSID")
	twilioAuthToken := os.Getenv("twilioAuthToken")
	twilioPhoneNumber := os.Getenv("twilioPhoneNumber")

	// Twilio API endpoint
	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", twilioAccountSID)

	// Prepare the form data
	formData := url.Values{}
	formData.Set("To", toPhoneNumber)
	formData.Set("From", twilioPhoneNumber)
	formData.Set("Body", messageBody)

	// fmt.Println("toPhoneNumber", toPhoneNumber)
	// fmt.Println("twilioPhoneNumber", twilioPhoneNumber)

	// Create the HTTP request
	req, err := http.NewRequest("POST", twilioURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create Twilio request: %v", err)
	}
	req.SetBasicAuth(twilioAccountSID, twilioAuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request to Twilio
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Twilio request: %v", err)
	}
	defer resp.Body.Close()

	// fmt.Println("err", err)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send WhatsApp notification via Twilio, status: %d", resp.StatusCode)
	}

	return nil
}
