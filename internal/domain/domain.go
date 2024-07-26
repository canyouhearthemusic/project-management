package domain

import (
	"fmt"
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Field   string `json:"field"`
}

type OnlyDate string

const DateLayout = "2006-01-02"

// method of [driver.Valuer] interface
func (o *OnlyDate) Value() (interface{}, error) {
	nt, err := time.Parse(DateLayout, string(*o))
	if err != nil {
		panic(err)
		// return nil, err
	}

	return nt, nil
}

// method of [sql.Scanner] interface
func (o *OnlyDate) Scan(val interface{}) error {
	nt, ok := val.(time.Time)
	if !ok {
		return fmt.Errorf("expected time.Time, got %T", val)
	}

	*o = OnlyDate(nt.Format(DateLayout))

	return nil
}

func (o *OnlyDate) String() string {
	return string(*o)
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("Field %s has issue: %s", e.Field, e.Message)
}
