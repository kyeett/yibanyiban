package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateIBAN(t *testing.T) {
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
func TestInvalidIBANError(t *testing.T) {
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

func TestIBANHandler(t *testing.T) {
	tcs := []struct {
		name string
		path string
	}{
		{"", "validate?iban=GB82WEST12345698765432"},
		{"", "validate?iban=AL85751639367318444714198669"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8080/"+tc.path, nil)
			if err != nil {
				t.Fatalf("failed to create http request, %v", err)
			}
			rec := httptest.NewRecorder()

			validateIBANHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response, %v", err)
			}
		})
	}

}

// Test iban?
// Test iban
