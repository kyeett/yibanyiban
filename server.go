package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN

//https://www.mobilefish.com/download/iban/random_generated_iban.txt

var (
	errCheckSumIncorrect = errors.New("checksum is incorrect")
	errNumberTooShort    = errors.New("IBAN is too short, expected > 4")
	errNumberTooLong     = errors.New("IBAN is too long, expected < 34")
	errInvalidCharacters = errors.New("Invalid characters, allowed are alphanumeric (A-Z, 0-9) and space (' ')")
)

var validChars = regexp.MustCompile("^[A-Z0-9 ]*$")
var mod97 = big.NewInt(97)

func validateIBAN(code string) (bool, error) {
	if len(code) < 4 {
		return false, errNumberTooShort
	}

	if len(code) > 34 {
		return false, errNumberTooLong
	}

	code = strings.ToUpper(code)

	// Check for non-alpha numerics or space
	if !validChars.Match([]byte(code)) {
		return false, errInvalidCharacters
	}

	// Remove spaces
	code = strings.Replace(code, " ", "", -1)

	// Rearrange
	code = code[4:] + code[0:4]

	// Convert individual letters to integers
	// A = 10, B = 11, ..., Z = 35
	var b strings.Builder
	for _, r := range code {
		if r < 'A' {
			b.WriteRune(r)
		} else {
			b.WriteString(strconv.Itoa(int(r - 'A' + 10)))
		}
	}

	// Verify checksum, mod 97 = 1
	numeric, _ := new(big.Int).SetString(b.String(), 10)
	remainder := uint(numeric.Mod(numeric, mod97).Uint64())
	if remainder != 1 {
		return false, errCheckSumIncorrect
	}

	return true, nil
}

type ValidationResponse struct {
	IBAN    string `json:"iban"`
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

func validateIBANHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	IBAN, ok := r.URL.Query()["iban"]
	if !ok {
		http.Error(w, "URL param 'iban' is missing", http.StatusForbidden)
		return
	}

	if len(IBAN) > 1 {
		http.Error(w, "multiple URL param 'iban' specified, only one is supported", http.StatusForbidden)
		return
	}

	valid, err := validateIBAN(IBAN[0])
	message := "OK"
	if err != nil {
		message = err.Error()
	}

	j, err := json.Marshal(ValidationResponse{
		IBAN:    IBAN[0],
		Valid:   valid,
		Message: message,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func main() {
	http.HandleFunc("/", validateIBANHandler)
	srv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
	}
	log.Fatal(srv.ListenAndServe())

	// country code using ISO 3166-1 alpha-2 – two letters,
	// check digits – two digits, and
	// Basic Bank Account Number (BBAN) – up to 30 alphanumeric characters that are country-specific.[1]
	// THe api converts upp
}
