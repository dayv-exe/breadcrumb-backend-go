package utils

import (
	"breadcrumb-backend-go/constants"
	"time"
)

func GetTimeNow() string {
	return time.Now().Format(constants.FULL_DATE_TIME_LAYOUT)
}

func ConvertToDate(s string) (time.Time, error) {
	return time.Parse(constants.DATE_LAYOUT, s)
}
