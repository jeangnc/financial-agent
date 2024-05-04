package btg

import (
	"testing"

	"github.com/jeangnc/financial-agent/pdf"
)

func TestParseFile(t *testing.T) {
	f := pdf.File{
		Pages: []pdf.Page{
			pdf.Page{Content: ""},
		},
	}
	ParseFile(f)
}
