package main

import (
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

// Test the HTTP status codes from the IBANHandler
func TestIBANHandler(t *testing.T) {
	tcs := []struct {
		name               string
		method             string
		path               string
		expectedStatusCode int
	}{
		{"ok #1", "GET", "validate?iban=GB82WEST12345698765432", http.StatusOK},
		{"ok #2", "GET", "validate?iban=AL85751639367318444714198669", http.StatusOK},
		{"invalid method", "POST", "validate?iban=AL85751639367318444714198669", http.StatusMethodNotAllowed},
		{"missing parameter iban", "GET", "validate?lolban=AL85751639367318444714198669", http.StatusForbidden},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "localhost:8080/"+tc.path, nil)
			rec := httptest.NewRecorder()

			validateIBANHandler(rec, req)
			res := rec.Result()

			if res.StatusCode != tc.expectedStatusCode {
				t.Fatalf("expected '%s', got '%s'\n", http.StatusText(tc.expectedStatusCode), http.StatusText(res.StatusCode))
			}
		})
	}

}
