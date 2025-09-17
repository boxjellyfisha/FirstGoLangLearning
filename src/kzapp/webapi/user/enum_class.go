package user

import (
	"encoding/json"
	"errors"
	"time"
)

type EnumClass string

func (e EnumClass) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return e.UnmarshalText([]byte(s))
}

func (e EnumClass) UnmarshalText(b []byte) error {
	if string(b) == "User" {
		e = EnumClass_User
	} else if string(b) == "Admin" {
		e = EnumClass_Admin
	} else {
		return errors.New("invalid enum value")
	}
	return nil
}

const (
	EnumClass_User  EnumClass = "User"
	EnumClass_Admin EnumClass = "Admin"
)

func (e EnumClass) String() string {
	return string(e)
}

type EnumData struct {
	Enum EnumClass `json:"enum"`
	Date time.Time `json:"date"`
}
