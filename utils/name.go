package utils

import (
	"breadcrumb-backend-go/constants"
	"strings"
	"time"
)

func NameIsValid(name *string) bool {
	return len(*name) == 0 || len(*name) <= constants.MAX_FULLNAME_CHARS
}

func NameChangeAllowed(lastChangedOn string) (bool, error) {
	// format is dd/mm/yyyy
	if strings.TrimSpace(lastChangedOn) == "" {
		return true, nil
	}

	lastChangeDate, err := time.Parse(constants.DATE_LAYOUT, lastChangedOn)
	if err != nil {
		return false, err
	}

	changeUnfreezeDate := lastChangeDate.AddDate(0, 0, constants.NAME_CHANGE_FREEZE_TIME)

	if time.Now().After(changeUnfreezeDate) {
		return true, nil
	}

	return false, nil
}
