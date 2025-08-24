package utils

import (
	"strconv"
	"time"
)

type JsonEpochTime time.Time

func (t JsonEpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}
