package btg

import (
	"reflect"
	"testing"

	"github.com/jeangnc/financial-agent/pdf"
)

func TestParseFile(t *testing.T) {
	tests := map[string]struct {
		input  string
		result []T
	}{
		"simple transactions": {
			input: "01 Jan Transaction 1 (1/2) R$ 10,10 01 Abr Transaction 2 R$ 0,98",
			result: []T{
				T{
					"date":        "01 Jan",
					"description": "Transaction 1 (1/2)",
					"amount":      "R$ 10,10",
				},
				T{
					"date":        "01 Abr",
					"description": "Transaction 2",
					"amount":      "R$ 0,98",
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
