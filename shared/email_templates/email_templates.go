package email_templates

import (
	"html/template"
	"log"
)

func GetEmailTemplate(file string) *template.Template {
	log.Print(file)
	// Parse the HTML template file
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
		return nil
	}

	return tmpl
}
