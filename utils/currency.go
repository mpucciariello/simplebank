package utils

const (
	USD = "USD"
	ARS = "ARS"
	EUR = "EUR"
)

// IsSupported returns true if the currency is supported
func IsSupported(currency string) bool {
	switch currency {
	case USD, ARS, EUR:
		return true
	}
	return false
}
