package btg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jeangnc/financial-agent/currency"
	"github.com/jeangnc/financial-agent/date"
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
			transaction, err := buildTransaction(match)
			if err != nil {
				return nil, err
			}

			result = append(result, *transaction)
		}
	}

	return result, nil
}

func buildTransaction(match regexp.RegexpMatch) (*types.Transaction, error) {
	amount, err := currency.ParseBrl(match["amount"])
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount: %s", err)
	}

	date, err := date.Parse(match["date"])
	if err != nil {
		return nil, fmt.Errorf("failed to convert date: %s", err)
	}

	description, currentInstallment, totalInstallments, err := extractInstallements(match["description"])
	if err != nil {
		return nil, err
	}

	t := types.Transaction{
		Description:        description,
		Amount:             amount,
		Date:               date,
		CurrentInstallment: currentInstallment,
		TotalInstallments:  totalInstallments,
	}

	return &t, nil
}

func extractInstallements(description string) (string, int64, int64, error) {
	installmentMatch, _ := regexp.Match(INSTALLMENT_REGEXP, description)

	if installmentMatch != nil {
		current, err := strconv.ParseInt(installmentMatch["current"], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("failed to convert current installment: %s", err)
		}

		total, err := strconv.ParseInt(installmentMatch["total"], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("failed to convert total installments: %s", err)
		}

		newDescription := strings.TrimSpace(regexp.Remove(INSTALLMENT_REGEXP, description))
		return newDescription, current, total, nil
	}

	return description, 1, 1, nil
}
