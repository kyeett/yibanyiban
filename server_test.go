package main

import "testing"

func Test_validateIBAN(t *testing.T) {
	tcs := []struct {
		IBAN     string
		expected bool
	}{
		{"GB82WEST12345698765432", true},
		{"AL86751639367318444714198669", true},
		{"AL87751639367318444714198669", false},
	}

	for _, tc := range tcs {
		valid, _ := validateIBAN(tc.IBAN)
		if valid != tc.expected {
			t.Errorf("expected validateIBAN(%q)=%t, got %t", tc.IBAN, tc.expected, valid)
		}
	}
}
func Test_invalidIBANError(t *testing.T) {
	tcs := []struct {
		IBAN        string
		expectedErr error
	}{
		{"GB8", errNumberTooShort},
		{"", errNumberTooShort},
		{"AL86751639367318444714198669AL86751639367318444714198669", errNumberTooLong},
		{"SE120312301023012301203012301230120301230210300123010230", errNumberTooLong},
		{"GB82WEST12345698765432&", errInvalidCharacters},
		{"GB82WEST12345698765432*", errInvalidCharacters},
		{"XX82WEST12345698765432", errCountryCodeInvalid},
		{"GB83WEST12345698765432", errCheckSumIncorrect},
		{"AL85751639367318444714198669", errCheckSumIncorrect},
		{"GB82WEST12345698765432", nil},
	}

	for _, tc := range tcs {
		_, err := validateIBAN(tc.IBAN)
		if err != tc.expectedErr {
			t.Errorf("expected error %q from validateIBAN(%q), got %q", tc.expectedErr, tc.IBAN, err)
		}
	}
}

// Test iban?
// Test iban
