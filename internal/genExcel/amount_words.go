package genexcel

import "strings"

var ones = []string{
	"",
	"One",
	"Two",
	"Three",
	"Four",
	"Five",
	"Six",
	"Seven",
	"Eight",
	"Nine",
	"Ten",
	"Eleven",
	"Twelve",
	"Thirteen",
	"Fourteen",
	"Fifteen",
	"Sixteen",
	"Seventeen",
	"Eighteen",
	"Nineteen",
}

var tens = []string{
	"",
	"",
	"Twenty",
	"Thirty",
	"Forty",
	"Fifty",
	"Sixty",
	"Seventy",
	"Eighty",
	"Ninety",
}

func twoDigitWords(n int) string {
	if n < 20 {
		return ones[n]
	}

	result := tens[n/10]

	if n%10 != 0 {
		result += " " + ones[n%10]
	}

	return result
}

func threeDigitWords(n int) string {
	if n < 100 {
		return twoDigitWords(n)
	}

	result := ones[n/100] + " Hundred"

	if n%100 != 0 {
		result += " " + twoDigitWords(n%100)
	}

	return result
}

// AmountInWords converts a rupee amount to words using the Indian numbering
// system: Crore, Lakh, Thousand, and Hundred.
func AmountInWords(amount float64) string {
	n := int64(amount)

	if n == 0 {
		return "Rupees Zero Only"
	}

	crore := n / 10_000_000
	n %= 10_000_000

	lakh := n / 1_00_000
	n %= 1_00_000

	thousand := n / 1_000
	n %= 1_000

	hundred := n

	var parts []string

	if crore > 0 {
		parts = append(parts, twoDigitWords(int(crore))+" Crore")
	}

	if lakh > 0 {
		parts = append(parts, twoDigitWords(int(lakh))+" Lakh")
	}

	if thousand > 0 {
		parts = append(parts, twoDigitWords(int(thousand))+" Thousand")
	}

	if hundred > 0 {
		parts = append(parts, threeDigitWords(int(hundred)))
	}

	return "Rupees " + strings.Join(parts, " ") + " Only"
}
