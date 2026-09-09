package generator

import "github.com/Dimensionexpert/payslip/internal/concurrency"

func CountConversionResults(
	results []concurrency.ConversionResult,
) (success int, failed int) {
	for _, result := range results {
		if result.Err != nil {
			failed++
		} else {
			success++
		}
	}

	return success, failed
}
