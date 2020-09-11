package time

import "time"

//GetTimeNow return the date now formatted in time zone
func GetTimeNow(format, timeZone string, onlyDate bool) (string, error) {
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