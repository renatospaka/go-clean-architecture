package utils

import (
	"time"
)

const (
	LayoutISO     = "2006-01-02"
	LayoutUS      = "January 2, 2006"
	LayoutShort   = "Jan-02-06"
	LayoutBR      = "January 02, 2006"
	LayoutBRShort = "01/Jan/2006"
)

func IsDateEqualToday(d time.Time) bool {
	today := time.Now()
	day := d.Day() == today.Day() 
	month := d.Month() == today.Month() 
	year := d.Year() == today.Year()
	
	//d.Equal(today) && 
	return day && month && year
}

func IsDateGreaterToday(d time.Time) bool {
	today := time.Now()	
	return d.After(today)
}

func IsDateLowerToday(d time.Time) bool {
	today := time.Now()
	return d.Before(today)
	
}

func isDateGreaterEqualToday(d time.Time) bool {
	if IsDateEqualToday(d) {
		return true
	}
	return IsDateGreaterToday(d)
}

func isDateLowerEqualToday(d time.Time) bool {
	if IsDateEqualToday(d) {
		return true
	}
	return IsDateLowerToday(d)
}
