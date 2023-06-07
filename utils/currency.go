package utils

const (
	USD = "USD"
	ARS = "ARS"
	EUR = "EUR"
)

// IsSupported returns true if the currency is supported
func IsSupported(c string) bool {
	switch c {
	case USD, ARS, EUR:
		return true
	}
	return false
}
