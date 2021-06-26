package utils

import "time"

func IsNilString(stuff interface{}) string {
	if stuff == nil {
		return ""
	}
	return stuff.(string)
}

func IsNilTime(stuff interface{}) time.Time {
	if stuff == nil {
		return time.Time{}
	}
	return stuff.(time.Time)
}