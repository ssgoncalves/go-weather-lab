package domain

import "regexp"

type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
}

func IsValidZipCode(zipCode string) bool {
	regex := regexp.MustCompile(`^\d{8}$`)
	return regex.MatchString(zipCode)
}
