package btg

const (
	DATE_REGEXP        = `(?<date>[0-9]{2} [a-zA-Z]{3})`
	AMOUNT_REGEXP      = `(?<amount>(?:- ?)?[A-Z]{1,}\$ (?:[0-9]+\.?)+\,[0-9]{2})`
	DESCRIPTION_REGEXP = `(?<description>([^R]|R[^$])*)`
	TRANSACTION_REGEXP = `(` + DATE_REGEXP + `\W+` + DESCRIPTION_REGEXP + `\W+` + AMOUNT_REGEXP + `)`
	INSTALLMENT_REGEXP = `\((?<current>\d+)\/(?<total>\d+)\)$`
)
