package btg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jeangnc/financial-agent/pdf"
	"github.com/jeangnc/financial-agent/regexp"
	"github.com/jeangnc/financial-agent/types"
)

const (
	DATE_REGEXP        = `(?<date>[0-9]{2} [a-zA-Z]{3})`
	AMOUNT_REGEXP      = `(?<amount>(?:- ?)?[A-Z]{1,}\$ (?:[0-9]+\.?)+\,[0-9]{2})`
	DESCRIPTION_REGEXP = `(?<description>([^R]|R[^$])*)`
	TRANSACTION_REGEXP = `(` + DATE_REGEXP + `\W+` + DESCRIPTION_REGEXP + `\W+` + AMOUNT_REGEXP + `)`
	INSTALLMENT_REGEXP = `\((?<current>\d+)\/(?<total>\d+)\)$`
)

func ParseFile(f pdf.File) ([]types.Transaction, error) {
	var result = make([]types.Transaction, 0)

	for _, page := range f.Pages {
		for _, match := range regexp.MatchAll(TRANSACTION_REGEXP, page.Content) {
			t, err := buildTransaction(match)
			if err != nil {
				return nil, fmt.Errorf("failed to convert amount: %s", err)
			}

			result = append(result, *t)
		}
	}

	return result, nil
}

func buildTransaction(match regexp.RegexpMatch) (*types.Transaction, error) {
	amount, err := parseCurrency(match["amount"])
	if err != nil {
		return nil, err
	}

	t := types.Transaction{
		Description:        match["description"],
		Amount:             amount,
		Date:               time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		CurrentInstallment: 1,
		TotalInstallments:  1,
	}

	installmentMatch, _ := regexp.Match(INSTALLMENT_REGEXP, t.Description)

	if installmentMatch != nil {
		currentInstallment, _ := strconv.ParseInt(installmentMatch["current"], 10, 64)
		totalInstallments, _ := strconv.ParseInt(installmentMatch["total"], 10, 64)

		t.Description = strings.TrimSpace(regexp.Remove(INSTALLMENT_REGEXP, t.Description))
		t.CurrentInstallment = currentInstallment
		t.TotalInstallments = totalInstallments
	}

	return &t, nil
}

func parseCurrency(amountStr string) (float64, error) {
	amountStr = strings.ReplaceAll(amountStr, "R$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	amountStr = strings.TrimSpace(amountStr)
	return strconv.ParseFloat(amountStr, 64)
}
