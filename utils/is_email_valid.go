package utils

import "net/mail"

func IsEmailValid(email string) bool {
	// returns true if email is valid
	_, err := mail.ParseAddress(email)
	return err == nil
}
