package types

import "time"

type Transaction struct {
	Date               time.Time
	Description        string
	Amount             float64
	CurrentInstallment int64
	TotalInstallments  int64
}
