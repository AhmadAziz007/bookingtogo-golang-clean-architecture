package model

import (
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(dateFormat) + `"`), nil
}

func (d Date) String() string {
	return d.Time.Format(dateFormat)
}
