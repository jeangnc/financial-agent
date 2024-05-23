package btg

import (
	"reflect"
	"testing"

	"github.com/jeangnc/financial-agent/pdf"
	"github.com/jeangnc/financial-agent/types"
)

func TestParseFile(t *testing.T) {
	tests := map[string]struct {
		input  string
		result []types.Transaction
	}{
		"simple transactions": {
			input: "01 Jan Transaction 1 (1/2) R$ 10,10 01 Abr Transaction 2 R$ 0,98",
			result: []types.Transaction{
				types.Transaction{
					"date":                "01 Jan",
					"description":         "Transaction 1",
					"amount":              "R$ 10,10",
					"current_installment": "1",
					"total_installments":  "2",
				},
				types.Transaction{
					"date":                "01 Abr",
					"description":         "Transaction 2",
					"amount":              "R$ 0,98",
					"current_installment": "1",
					"total_installments":  "1",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := pdf.File{
				Pages: []pdf.Page{
					pdf.Page{Content: test.input},
				},
			}
			transactions := ParseFile(f)

			if !reflect.DeepEqual(transactions, test.result) {
				t.Fatalf("returned %q; expected %q", transactions, test.result)
			}
		})
	}
}
