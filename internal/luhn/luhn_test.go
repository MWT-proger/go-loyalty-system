package luhn

import (
	"testing"
)

func TestValidate(t *testing.T) {
	validNumber := "79927398713"
	invalidNumber := "12345678901"

	valid := Validate(validNumber)
	if !valid {
		t.Errorf("expected %s to be valid", validNumber)
	}

	invalid := Validate(invalidNumber)
	if invalid {
		t.Errorf("expected %s to be invalid", invalidNumber)
	}
}

func TestCalculateLuhnSum(t *testing.T) {
	number := "79927398713"
	parity := len(number) % 2

	expectedSum := int64(70)
	sum, err := calculateLuhnSum(number, parity)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sum != expectedSum {
		t.Errorf("expected sum to be %d, got %d", expectedSum, sum)
	}
}
