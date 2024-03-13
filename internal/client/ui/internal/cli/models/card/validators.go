package card

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (m Model) validateFields() error {
	if len(m.inputs[name].Value()) == 0 {
		return fmt.Errorf("empty name")
	}

	if err := ccnValidator(m.inputs[ccn].Value()); err != nil {
		return fmt.Errorf("wrong ccn field: %w", err)
	}

	if err := expValidator(m.inputs[exp].Value()); err != nil {
		return fmt.Errorf("wrong exp field: %w", err)
	}

	if err := cvvValidator(m.inputs[cvv].Value()); err != nil {
		return fmt.Errorf("wrong cvv field: %w", err)
	}

	return nil
}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func getCcn(s string) string {
	return s
}

func getExp(s string) time.Time {
	t, _ := time.Parse("01/06", s)
	return t
}

func getCvv(s string) uint32 {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint32(v)
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("EXP is invalid")
	}

	// There should be only one slash and it should be in the 2nd index (3rd character)
	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("EXP is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	// The CVV should be a number of 3 digits
	// Since the input will already ensure that the CVV is a string of length 3,
	// All we need to do is check that it is a number
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}
