package util

const (
	EUR = "EUR"
	CAD = "CAD"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}

	return false
}
