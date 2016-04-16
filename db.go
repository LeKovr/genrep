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
	Id       string  `json:"id"`
	Customer string  `json:"customer"`
	Date     Date    `json:"date"`
	Tax      float64 `json:"tax"`
	Services []struct {
		Amount float64 `json:"Amount"`
		Code   string  `json:"Code"`
	} `json:"services"`
}

type Party struct {
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
		Author   string `json:"Author"`
		Creator  string `json:"Creator"`
		Font     string `json:"Font"`
		Keywords string `json:"Keywords"`
		Subject  string `json:"Subject"`
		Title    string `json:"Title"`
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

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	d.Time, err = time.Parse(dtLayout, string(b))
	return
}
