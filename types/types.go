package types

import "time"

type Transaction struct {
	Date               time.Time
	Description        string
	Amount             float64
	CurrentInstallment int64
	TotalInstallments  int64
}

type Category struct {
	Name string
}

type ExternalTransaction struct {
	Id                 int64
	Date               time.Time
	Description        string
	Amount             float64
	CurrentInstallment int64
	TotalInstallments  int64
	Category           Category
}
