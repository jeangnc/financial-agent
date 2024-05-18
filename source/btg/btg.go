package btg

import (
	"fmt"
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
		for _, m := range matchAll(TRANSACTION_REGEXP, p.Content) {
			m2, _ := match(INSTALLMENT_REGEXP, m["description"])

			if m2 != nil {
				m["current_installment"] = m2["current"]
				m["total_installments"] = m2["total"]
			}

			result = append(result, m)
		}
	}

	return result
}

func match(pattern string, text string) (T, error) {
	matches := matchAll(pattern, text)

	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches for the pattern: %s", pattern)
	}

	return matches[0], nil
}

func matchAll(pattern string, text string) []T {
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
