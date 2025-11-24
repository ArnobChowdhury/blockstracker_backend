package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// JSONTime is a custom time.Time wrapper with a critical purpose: to enforce
// consistent timestamp formatting across the entire application stack.
//
// ### Why does this exist? ###
// The client-side applications (React Native, Electron) use JavaScript's
// `Date.toISOString()`, which always includes millisecond precision
// (e.g., "2025-11-20T18:00:00.000Z").
//
// Go's default JSON marshaling for `time.Time` omits zero-value milliseconds
// (e.g., "2025-11-20T18:00:00Z").
//
// This subtle difference caused a critical bug where the UNIQUE constraint on
// tasks (`repetitive_task_template_id`, `due_date`) failed during sync because
// the string representations of the timestamps did not match.
//
// This type ensures all timestamps sent from the backend include milliseconds,
// aligning with the clients and preventing sync errors.
//
// DO NOT REMOVE THIS TYPE or replace it with the standard time.Time in models
// without understanding these consequences. See ADR-001 for more details.
type JSONTime time.Time

const jsonTimeFormat = "2006-01-02T15:04:05.000Z"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}

	stamp := fmt.Sprintf("\"%s\"", time.Time(t).UTC().Format(jsonTimeFormat))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	*t = JSONTime(parsedTime)
	return nil
}

func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *JSONTime) Scan(value any) error {
	if value == nil {
		*t = JSONTime(time.Time{})
		return nil
	}
	if vt, ok := value.(time.Time); ok {
		*t = JSONTime(vt)
		return nil
	}
	return fmt.Errorf("failed to scan JSONTime: value is not time.Time")
}
