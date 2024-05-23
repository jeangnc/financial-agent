package btg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	r "regexp"

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
				return nil, fmt.Errorf("failed to convert amount: %s", err)
			}

			result = append(result, *transaction)
		}
	}

	return result, nil
}

func buildTransaction(match regexp.RegexpMatch) (*types.Transaction, error) {
	amount, err := parseBrlCurrency(match["amount"])
	if err != nil {
		return nil, err
	}

	date, err := parseBrlDate(match["date"])
	if err != nil {
		return nil, err
	}

	description, currentInstallment, totalInstallments := extractInstallements(match["description"])
	t := types.Transaction{
		Description:        description,
		Amount:             amount,
		Date:               date,
		CurrentInstallment: currentInstallment,
		TotalInstallments:  totalInstallments,
	}

	return &t, nil
}

func extractInstallements(description string) (string, int64, int64) {
	installmentMatch, _ := regexp.Match(INSTALLMENT_REGEXP, description)

	if installmentMatch != nil {
		current, _ := strconv.ParseInt(installmentMatch["current"], 10, 64)
		total, _ := strconv.ParseInt(installmentMatch["total"], 10, 64)

		newDescription := strings.TrimSpace(regexp.Remove(INSTALLMENT_REGEXP, description))
		return newDescription, current, total
	}

	return description, 1, 1
}

func parseBrlDate(dateStr string) (time.Time, error) {
	months := map[string]int{
		"Jan": 1,
		"Fev": 2,
		"Mar": 3,
		"Abr": 4,
		"Mai": 5,
		"Jun": 6,
		"Jul": 7,
		"Ago": 8,
		"Set": 9,
		"Out": 10,
		"Nov": 11,
		"Dez": 12,
	}

	var nonAlpha = r.MustCompile(`\d`)
	monthString := strings.TrimSpace(nonAlpha.ReplaceAllString(dateStr, ""))
	month, ok := months[monthString]
	if !ok {
		return time.Now(), fmt.Errorf("failed to convert month: %s", monthString)
	}

	var nonNumeric = r.MustCompile(`[a-zA-Z ]`)
	dayString := strings.TrimSpace(nonNumeric.ReplaceAllString(dateStr, ""))
	day, err := strconv.ParseInt(dayString, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(2024, time.Month(month), int(day), 0, 0, 0, 0, time.Local), nil
}

func parseBrlCurrency(amountStr string) (float64, error) {
	amountStr = strings.ReplaceAll(amountStr, "R$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	amountStr = strings.TrimSpace(amountStr)
	return strconv.ParseFloat(amountStr, 64)
}
