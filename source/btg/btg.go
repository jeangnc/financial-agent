package btg

import (
	"regexp"

	"github.com/jeangnc/financial-agent/pdf"
)

const DATE_REGEXP = `(?<date>[0-9]{2} [a-zA-Z]{3})`
const AMOUNT_REGEXP = `(?<amount>(?:- ?)?[A-Z]{1,}\$ (?:[0-9]+\.?)+\,[0-9]{2})`
const DESCRIPTION_REGEXP = `(?<description>([^R]|R[^$])*)`
const TRANSACTION_REGEXP = `(` + DATE_REGEXP + `\W+` + DESCRIPTION_REGEXP + `\W+` + AMOUNT_REGEXP + `)`
const INSTALLMENT_REGEXP = `\((?<current>\d+)\/(?<total>\d+)\)$`

type T map[string]string

func ParseFile(f pdf.File) []T {
	var result = make([]T, 0)

	for _, p := range f.Pages {
		for _, m := range runRegexp(TRANSACTION_REGEXP, p.Content) {
			result = append(result, m)
		}
	}

	return result
}

func runRegexp(pattern string, text string) []T {
	result := make([]T, 0)

	var expr = regexp.MustCompile(pattern)
	for _, m := range expr.FindAllStringSubmatch(text, -1) {
		t := T{}

		for i, name := range expr.SubexpNames() {
			if name != "" {
				t[name] = m[i]
			}
		}

		result = append(result, t)
	}

	return result
}
