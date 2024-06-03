package organizze

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type (
	Date struct {
		time.Time
	}

	Category struct {
		Id       int64  `json:"id"`
		Archived bool   `json:"archived"`
		ParentID int64  `json:"parent_id"`
		Name     string `json:"name"`
	}

	InstallmentAttributes struct {
		Periodicity string  `json:"periodicity"`
		Total       float64 `json:"total"`
	}

	TagString string

	Tag struct {
		Name string `json:"name"`
	}

	Transaction struct {
		Id                    int64                  `json:"id,omitempty"`
		CreditCardId          int64                  `json:"credit_card_id"`
		CreditCardInvoiceId   int64                  `json:"credit_card_invoice_id"`
		Description           string                 `json:"description"`
		Notes                 string                 `json:"notes"`
		Date                  Date                   `json:"date"`
		AmountCents           float64                `json:"amount_cents"`
		CategoryId            int64                  `json:"category_id"`
		InstallmentAttributes *InstallmentAttributes `json:"installments_attributes"`
		Tags                  TagString              `json:"tags"`
	}

	InvoiceHeader struct {
		Id           int64 `json:"id",omitempty`
		Date         Date  `json:"date"`
		StartingDate Date  `json:"starting_date"`
		ClosingDate  Date  `json:"closing_date"`
	}

	Invoice struct {
		Id           int64         `json:"id",omitempty`
		Date         Date          `json:"date"`
		StartingDate Date          `json:"starting_date"`
		ClosingDate  Date          `json:"closing_date"`
		Transactions []Transaction `json:"transactions"`
	}
)

func (t TagString) MarshalJSON() ([]byte, error) {
	tags := make([]Tag, 0, 0)
	splits := strings.Split(string(t), ",")

	for _, s := range splits {
		tags = append(tags, Tag{Name: s})
	}
	return json.Marshal(tags)
}

/*(
func (t *TagString) Tags() []Tag {
	tags := make([]Tag, 0, 0)
	splits := strings.Split(string(*t), ",")

	for _, s := range splits {
		tags = append(tags, Tag{Name: s})
	}

	return tags
}
*/

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	d.Time, err = time.Parse(DateLayout, s)
	return
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format(DateLayout))), nil
}
