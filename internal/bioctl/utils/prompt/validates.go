package prompt

import (
	"errors"
	"strconv"
)

func OptionalInput(ans interface{}) error {
	return nil
}

func ValidateIntegerNumberInput(ans interface{}) error {
	_, err := strconv.ParseInt(ans.(string), 0, 64)
	if err != nil {
		return errors.New("invalid number")
	}
	return nil
}
