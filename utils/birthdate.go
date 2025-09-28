package utils

import (
	"breadcrumb-backend-go/constants"
	"time"
)

func BirthdateIsValid(date string) (bool, error) {
	// returns true if date is valid
	// format is dd/mm/yyyy
	birthdate, err := time.Parse(constants.BIRTHDATE_ONLY_DATE_LAYOUT, date)

	if err != nil {
		return false, err
	}

	minAgeDate := time.Now().AddDate(-constants.MIN_AGE, 0, 0)
	maxAgeDate := time.Now().AddDate(-constants.MAX_AGE, 0, 0)

	if birthdate.After(maxAgeDate) || birthdate.Equal(maxAgeDate) {
		if birthdate.Before(minAgeDate) || birthdate.Equal(minAgeDate) {
			return true, nil
		}
	}

	return false, nil
}
