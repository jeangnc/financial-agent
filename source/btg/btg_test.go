package btg

import (
	"reflect"
	"testing"
	"time"

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
					Date:               time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
					Description:        "Transaction 1",
					Amount:             10.10,
					CurrentInstallment: 1,
					TotalInstallments:  2,
				},
				types.Transaction{
					Date:               time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
					Description:        "Transaction 2",
					Amount:             0.98,
					CurrentInstallment: 1,
					TotalInstallments:  1,
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
			transactions, _ := ParseFile(f)

			if !reflect.DeepEqual(transactions, test.result) {
				t.Fatalf("returned %#v; expected %#v", transactions, test.result)
			}
		})
	}
}
