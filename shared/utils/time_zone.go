package util

import (
	"log"
	"time"
)

func GetJakartaTimeZone() (*time.Location, error) {
	// Load the GMT+7 location
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return nil, err
	}

	return location, nil
}
