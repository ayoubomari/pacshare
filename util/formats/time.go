package formats

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// convert from string:"00:04:50" to Number:290
func DurationStrToSeconds(durationStr string) int {
	if durationStr == "" {
		return 0
	}

	durationArray := strings.Split(durationStr, ":")

	s := 0
	p := 1
	for i := len(durationArray) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(durationArray[i])
		s += n * p
		p *= 60
	}

	return s
}

// convert from number:6 to string:"06"
func Ecritetwo(num int) string {
	numberString := fmt.Sprint(num)
	if len(numberString) < 2 {
		return "0" + numberString
	} else {
		return numberString
	}
}

// convert from Number:290 to string:"00:04:50"
func DisplaySecends(sec int) string {
	var hours, mins int
	mins = sec / 60
	sec %= 60
	if mins >= 60 {
		hours = mins / 60
		mins %= 60
	} else {
		hours = 0
	}

	hoursString := Ecritetwo(hours)
	minsString := Ecritetwo(mins)
	secString := Ecritetwo(sec)

	return fmt.Sprintf("%s:%s:%s", hoursString, minsString, secString)
}

// convert timeStamp to human readable date
func TimeStampToDate(unixTimestamp int64) string {
	// Convert Unix timestamp to time.Time
	timeValue := time.Unix(unixTimestamp, 0)

	// Format the time in the desired layout
	// "Monday, January 2, 2006 15:04:05 PM"
	return timeValue.Format("Monday, January 2, 2006 15:04:05 PM")
}
