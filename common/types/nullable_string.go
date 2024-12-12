package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NullableString struct {
	sql.NullString
}

// NewNullableString creates a NullableString from a string pointer
func NewNullableString(s *string) NullableString {
	if s == nil {
		return NullableString{
			sql.NullString{
				String: "",
				Valid:  false,
			},
		}
	}
	return NullableString{
		sql.NullString{
			String: *s,
			Valid:  true,
		},
	}
}

// FromString creates a NullableString from a string value
func FromString(s string) NullableString {
	return NullableString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

// SetValid sets both the string value and valid flag
func (ns *NullableString) SetValid(s string) {
	ns.NullString = sql.NullString{
		String: s,
		Valid:  true,
	}
}

// SetInvalid marks the string as invalid
func (ns *NullableString) SetInvalid() {
	ns.NullString = sql.NullString{
		String: "",
		Valid:  false,
	}
}

// GetValue returns the string value or empty string if not valid
func (ns NullableString) GetValue() string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// MarshalJSON implements json.Marshaler interface
func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal("")
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ns *NullableString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		ns.SetInvalid()
		return nil
	}
	ns.SetValid(*s)
	return nil
}

// Value implements the driver.Valuer interface
func (ns NullableString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// Scan implements the sql.Scanner interface
func (ns *NullableString) Scan(value interface{}) error {
	return ns.NullString.Scan(value)
}
