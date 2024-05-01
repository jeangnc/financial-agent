package pdf

import (
	"bytes"
	"os"

	"github.com/dslipak/pdf"
)

type File struct {
	reader *pdf.Reader
}

type Page struct {
	Content string
}

func NewEncryptedReader(path string, password string) (*File, error) {
	r, err := openEncrypted(path, password)
	if err != nil {
		return nil, err
	}

	return &File{
		reader: r,
	}, nil
}

func openEncrypted(path string, password string) (*pdf.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}

	pw := func() string {
		return password
	}
	r, err := pdf.NewReaderEncrypted(f, fs.Size(), pw)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (pdf File) ReadPages() ([]Page, error) {
	pages := make([]Page, 0)

	for i := 0; i < pdf.reader.NumPage(); i++ {
		p, err := pdf.ReadPage(i)
		if err != nil {
			return nil, err
		}

		if p.Content == "" {
			continue
		}

		pages = append(pages, *p)
	}

	return pages, nil
}

func (pdf File) ReadPage(i int) (*Page, error) {
	p := pdf.reader.Page(i + 1)
	if p.V.IsNull() {
		return &Page{Content: ""}, nil
	}

	var buf bytes.Buffer
	rows, err := p.GetTextByRow()
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		buf.WriteString("\n")
		for _, word := range row.Content {
			buf.WriteString(word.S + " ")
		}
	}

	return &Page{Content: buf.String()}, nil
}
