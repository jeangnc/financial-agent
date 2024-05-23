package btg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jeangnc/financial-agent/pdf"
	"github.com/jeangnc/financial-agent/types"
)

const DATE_REGEXP = `(?<date>[0-9]{2} [a-zA-Z]{3})`
const AMOUNT_REGEXP = `(?<amount>(?:- ?)?[A-Z]{1,}\$ (?:[0-9]+\.?)+\,[0-9]{2})`
const DESCRIPTION_REGEXP = `(?<description>([^R]|R[^$])*)`
const TRANSACTION_REGEXP = `(` + DATE_REGEXP + `\W+` + DESCRIPTION_REGEXP + `\W+` + AMOUNT_REGEXP + `)`
const INSTALLMENT_REGEXP = `\((?<current>\d+)\/(?<total>\d+)\)$`

func ParseFile(f pdf.File) []types.Transaction {
	var result = make([]types.Transaction, 0)

	for _, p := range f.Pages {
		expr := regexp.MustCompile(TRANSACTION_REGEXP)
		for _, m := range matchAll(expr, p.Content) {
			expr = regexp.MustCompile(INSTALLMENT_REGEXP)
			m2, _ := match(expr, m["description"])

			if m2 != nil {
				m["description"] = strings.Trim(expr.ReplaceAllString(m["description"], ""), " ")
				m["current_installment"] = m2["current"]
				m["total_installments"] = m2["total"]
			} else {
				m["current_installment"] = "1"
				m["total_installments"] = "1"
			}

			result = append(result, m)
		}
	}

	return result
}

func match(expr *regexp.Regexp, text string) (types.Transaction, error) {
	matches := matchAll(expr, text)

	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches for the pattern: %s", expr)
	}

	return matches[0], nil
}

func matchAll(expr *regexp.Regexp, text string) []types.Transaction {
	result := make([]types.Transaction, 0)

	for _, m := range expr.FindAllStringSubmatch(text, -1) {
		t := types.Transaction{}

		for i, name := range expr.SubexpNames() {
			if name != "" {
				t[name] = m[i]
			}
		}

		result = append(result, t)
	}

	return result
}
