package reports

import (
	"time"
)

type Date struct {
	time.Time
}

const dtLayout = "2006-01-02"

type ServiceDef struct {
	Name  string  `json:"name"`
	Units string  `json:"units"`
	Price float64 `json:"price"`
}
type Document struct {
	Id          string  `json:"id"`
	Customer    string  `json:"customer"`
	Date        Date    `json:"date"`
	Tax         float64 `json:"tax"`
	ConfirmDate Date    `json:"confirmed"`
	Services    []struct {
		Amount float64 `json:"Amount"`
		Code   string  `json:"Code"`
	} `json:"services"`
}

type Party struct {
	Id       string `json:"Id"`
	Address  string `json:"Address"`
	Bank     string `json:"Bank"`
	Bik      string `json:"Bik"`
	Inn      string `json:"Inn"`
	Kpp      string `json:"Kpp"`
	Ks       string `json:"Ks"`
	Name     string `json:"Name"`
	NameFull string `json:"NameFull"`
	NameLong string `json:"NameLong"`
	Phone    string `json:"Phone"`
	Rs       string `json:"Rs"`
	Title    string `json:"Title"`
}

type CustomerServices map[string]ServiceDef

type Act struct {
	Contractor Party               `json:"contractor"`
	Customers  map[string]Party    `json:"customers"`
	Documents  map[string]Document `json:"documents"`
	Meta       struct {
		AuthorFmt  string `json:"AuthorFmt"`
		SubjectFmt string `json:"SubjectFmt"`
		TitleFmt   string `json:"TitleFmt"`
		Creator    string `json:"Creator"`
		Font       string `json:"Font"`
		Keywords   string `json:"Keywords"`
	} `json:"meta"`
	Names struct {
		Resume        string `json:"Resume"`
		Sign          string `json:"Sign"`
		Stamp         string `json:"Stamp"`
		TitleFmt      string `json:"TitleFmt"`
		TotalTax      string `json:"TotalTax"`
		TotalTitle    string `json:"TotalTitle"`
		TotalWordsFmt string `json:"TotalWordsFmt"`
		ZeroTax       string `json:"ZeroTax"`
	} `json:"names"`
	Services map[string]CustomerServices
}

type ActDef struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Customer    string  `json:"customer"`
	Contractor  string  `json:"contractor"`
	Date        Date    `json:"date"`
	ConfirmDate Date    `json:"confirm"`
	Amount      float64 `json:"amount"`
	Tax         float64 `json:"tax"`
}

// -----------------------------------------------------------------------------

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	d.Time, err = time.Parse(dtLayout, string(b))
	return
}

func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 12)
	b = append(b, '"')
	b = d.AppendFormat(b, "2006-01-02")
	b = append(b, '"')
	return b, nil
}
