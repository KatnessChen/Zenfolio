package utils

import "regexp"

var (
	SymbolRegex   = regexp.MustCompile(`^[A-Z0-9]{1,8}(\.[A-Z0-9]{1,2})?$`)
	CurrencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)
)
