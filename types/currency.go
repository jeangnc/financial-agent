package types

import (
	"strconv"
	"strings"
)

func ParseBrl(amountStr string) (float64, error) {
	amountStr = strings.ReplaceAll(amountStr, "R$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	amountStr = strings.TrimSpace(amountStr)
	return strconv.ParseFloat(amountStr, 64)
}
