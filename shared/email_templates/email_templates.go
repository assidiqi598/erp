package email_templates

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func GetEmailTemplateAndReplace(url string, data struct{}) (string, error) {
	// Fetch the template from the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch template: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch template: %s", resp.Status)
	}

	// Read the template content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Parse the template with placeholders
	tmpl, err := template.New("email").Parse(buf.String())
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Replace placeholders with actual values
	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}
