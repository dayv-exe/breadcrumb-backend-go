package utils

import "time"

func BirthdateIsValid(date string) (bool, error) {
	// format is dd/mm/yyyy
	layout := "02/01/2006"
	maxAge := 85
	minAge := 18
	birthdate, err := time.Parse(layout, date)

	if err != nil {
		return false, err
	}

	minAgeDate := time.Now().AddDate(-minAge, 0, 0)
	maxAgeDate := time.Now().AddDate(-maxAge, 0, 0)

	if birthdate.After(maxAgeDate) || birthdate.Equal(maxAgeDate) {
		if birthdate.Before(minAgeDate) || birthdate.Equal(minAgeDate) {
			return true, nil
		}
	}

	return false, nil
}
