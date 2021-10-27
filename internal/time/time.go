package time

import (
	"strconv"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(0, millis*int64(time.Millisecond)))

	return nil
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
