package types

import (
	"encoding/json"
	"strconv"
	"time"
)

// UnixString is a time.Time disguise as unix timestamp.
type UnixString time.Time

// Time returns time.Time.
func (s UnixString) Time() time.Time {
	return time.Time(s)
}

// Sec returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC.
func (s UnixString) Sec() int64 {
	return s.Time().Unix()
}

// String returns a unix timestamp string.
func (s UnixString) String() string {
	return strconv.FormatInt(s.Sec(), 10)
}

// MarshalJSON satisfies json.Marshaler interface.
func (s UnixString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON satisfies json.Unmarshaler interface.
func (s *UnixString) UnmarshalJSON(src []byte) error {
	var str string

	if err := json.Unmarshal(src, &str); err != nil {
		return err
	}

	ux, err := ParseUnixString(str)
	if err != nil {
		return err
	}

	*s = ux

	return nil
}

// UnixStringDate creates a new UnixString from a date time.
func UnixStringDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) UnixString {
	return UnixString(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

// UnixStringSec creates a new UnixString from a unix timestamp.
func UnixStringSec(sec int64) UnixString {
	return UnixString(time.Unix(sec, 0).UTC())
}

// ParseUnixString parses a unix timestamp string to UnixString.
func ParseUnixString(s string) (UnixString, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return UnixString{}, err
	}

	return UnixStringSec(sec), nil
}
