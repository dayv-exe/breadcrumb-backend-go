package utils

import "breadcrumb-backend-go/constants"

func NameIsValid(name *string) bool {
	return len(*name) == 0 || len(*name) <= constants.MAX_FULLNAME_CHARS
}
