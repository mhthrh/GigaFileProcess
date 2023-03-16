package Validation

import (
	"fmt"
	"github.com/almerlucke/go-iban/iban"
	"regexp"
)

var (
	maxAmount       float32 = 99999999
	isValidFullName         = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidaID(id interface{}) (int64, error) {
	i, ok := id.(int64)
	if ok {
		return 0, fmt.Errorf("cannot convert id to int64")
	}
	return i, nil
}
func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidateAmount(a interface{}) (float32, error) {
	amount, ok := a.(float32)
	if !ok {
		return 0, fmt.Errorf("amount value incorect")
	}

	if amount <= 0 {
		return 0, fmt.Errorf("amount is less than ziro/ is ziro")
	}
	if amount > maxAmount {
		return 0, fmt.Errorf("amount overflow")
	}
	return amount, nil
}

func ValidateIban(s string) error {
	_, err := iban.NewIBAN(s)

	if err != nil {
		return fmt.Errorf("check iban feild, %v", err)
	}
	return nil
}
