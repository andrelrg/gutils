package time

import "time"

//GetTimeNowString return the date now formatted in time zone
func GetTimeNowString(format, timeZone string, onlyDate bool) (string, error) {
	timeNow := time.Now()

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return "", err
	}
	timeNow = timeNow.In(loc)

	if onlyDate {
		timeNow = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())
	}

	timeString := timeNow.Format(format)
	return timeString, nil
}

func GetTimeNow(format, timeZone string) (time.Time, error) {
	timeNow := time.Now()

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, err
	}
	timeNow = timeNow.In(loc)

	return timeNow, nil
}