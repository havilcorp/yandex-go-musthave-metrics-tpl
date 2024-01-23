package errors

import (
	"fmt"
	"time"
)

type TimeError struct {
	Place string
	Time  time.Time
	Err   error
}

func (te *TimeError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

func NewTimeError(place string, err error) error {
	return &TimeError{
		Place: place,
		Time:  time.Now(),
		Err:   err,
	}
}
