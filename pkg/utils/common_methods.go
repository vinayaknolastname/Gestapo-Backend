package utils

import "time"

func ParseDate(dob string) (*time.Time, error) {
	if dob == "" {
		return nil, nil
	}
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return nil, err
	}
	return &parsedDOB, nil
}
