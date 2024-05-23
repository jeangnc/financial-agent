package btg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jeangnc/financial-agent/pdf"
	"github.com/jeangnc/financial-agent/types"
)

const DATE_REGEXP = `(?<date>[0-9]{2} [a-zA-Z]{3})`
const AMOUNT_REGEXP = `(?<amount>(?:- ?)?[A-Z]{1,}\$ (?:[0-9]+\.?)+\,[0-9]{2})`
const DESCRIPTION_REGEXP = `(?<description>([^R]|R[^$])*)`
const TRANSACTION_REGEXP = `(` + DATE_REGEXP + `\W+` + DESCRIPTION_REGEXP + `\W+` + AMOUNT_REGEXP + `)`
const INSTALLMENT_REGEXP = `\((?<current>\d+)\/(?<total>\d+)\)$`

type RegexpMatch map[string]string

func ParseFile(f pdf.File) ([]types.Transaction, error) {
	var result = make([]types.Transaction, 0)

	for _, p := range f.Pages {
		expr := regexp.MustCompile(TRANSACTION_REGEXP)
		for _, m := range matchAll(expr, p.Content) {
			amount, err := parseCurrency(m["amount"])
			if err != nil {
				return nil, fmt.Errorf("failed to convert amount: %s", err)
			}

			t := types.Transaction{
				Description:        m["description"],
				Amount:             amount,
				Date:               time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
				CurrentInstallment: 1,
				TotalInstallments:  1,
			}

			expr = regexp.MustCompile(INSTALLMENT_REGEXP)
			m2, _ := match(expr, t.Description)

			if m2 != nil {
				currentInstallment, _ := strconv.ParseInt(m2["current"], 10, 64)
				totalInstallments, _ := strconv.ParseInt(m2["total"], 10, 64)

				t.Description = strings.TrimSpace(expr.ReplaceAllString(t.Description, ""))
				t.CurrentInstallment = currentInstallment
				t.TotalInstallments = totalInstallments
			}

			result = append(result, t)
		}
	}

	return result, nil
}

func parseCurrency(amountStr string) (float64, error) {
	amountStr = strings.ReplaceAll(amountStr, "R$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	amountStr = strings.TrimSpace(amountStr)
	return strconv.ParseFloat(amountStr, 64)
}

func match(expr *regexp.Regexp, text string) (RegexpMatch, error) {
	matches := matchAll(expr, text)

	if len(matches) == 0 {
		return nil, nil
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple matches for the pattern: %s", expr)
	}

	return matches[0], nil
}

func matchAll(expr *regexp.Regexp, text string) []RegexpMatch {
	result := make([]RegexpMatch, 0)

	for _, m := range expr.FindAllStringSubmatch(text, -1) {
		t := RegexpMatch{}

		for i, name := range expr.SubexpNames() {
			if name != "" {
				t[name] = m[i]
			}
		}

		result = append(result, t)
	}

	return result
}
