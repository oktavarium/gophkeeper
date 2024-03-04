package cli

import (
	"fmt"
)

func ValidateInputs(inputs ...string) error {
	for _, input := range inputs {
		if len(input) == 0 {
			return fmt.Errorf("empty fields")
		}
	}
	return nil
}

func ValidatePasswords(first, second string) error {
	if first != second {
		return fmt.Errorf("passwords are not equal")
	}
	return nil
}
