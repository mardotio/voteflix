package utils

import (
	"strconv"
	"time"
)

type JsonEpochTime time.Time

func (t JsonEpochTime) TimestampString() string {
	return strconv.FormatInt(time.Time(t).Unix(), 10)
}

func (t JsonEpochTime) MarshalJSON() ([]byte, error) {
	return []byte(t.TimestampString()), nil
}
