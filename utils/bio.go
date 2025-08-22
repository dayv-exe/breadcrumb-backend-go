package utils

import "breadcrumb-backend-go/constants"

func BioIsValid(bio *string) bool {
	return len(*bio) <= constants.MAX_BIO_CHARS
}
