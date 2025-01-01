package storage

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"reflect"

	"github.com/minio/minio-go/v7"
)

func (s *S3ClientType) GetEmailTemplateAndReplace(bucketName string, objectKey string, data any) (string, error) {

	// Check if the provided data is a struct
	valData := reflect.ValueOf(data)
	if valData.Kind() != reflect.Struct {
		return "", fmt.Errorf("input is not a struct, got: %s", valData.Kind())
	}

	// Fetch the HTML file
	object, err := s.Client.GetObject(context.Background(), bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to get object: %v", err)
	}
	defer object.Close()

	// Read the template content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(object)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}
	htmlContent := buf.String()

	// Parse the template and replace placeholders
	tmpl, err := template.New(objectKey).Parse(htmlContent)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	var renderedHTML bytes.Buffer
	err = tmpl.Execute(&renderedHTML, data)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return htmlContent, nil
}
