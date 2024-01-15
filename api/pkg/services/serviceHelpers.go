package services

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordStrength(password string) bool {
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	numberRegex := regexp.MustCompile(`[0-9]`)
	symbolRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)

	strong := lowercaseRegex.MatchString(password) && uppercaseRegex.MatchString(password) && numberRegex.MatchString(password) && symbolRegex.MatchString(password)

	return strong
}
