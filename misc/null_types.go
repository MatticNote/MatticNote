package misc

import (
	"database/sql"
	"encoding/json"
	"time"
)

type (
	NullString struct {
		sql.NullString
	}
	NullTime struct {
		sql.NullTime
	}
	NullBool struct {
		sql.NullBool
	}
)

// NullString

func (ns *NullString) UnmarshalJSON(value []byte) error {
	err := json.Unmarshal(value, &ns.String)
	ns.Valid = err == nil
	return err
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

// NullTime

func (nt *NullTime) UnmarshalJSON(value []byte) error {
	err := json.Unmarshal(value, &nt.Time)
	nt.Valid = err == nil
	return err
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time.Format(time.RFC3339))
}

// NullBool

func (nb *NullBool) UnmarshalJSON(value []byte) error {
	err := json.Unmarshal(value, &nb.Bool)
	nb.Valid = err == nil
	return err
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nb.Bool)
}
