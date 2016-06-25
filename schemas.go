// Схемы данных пакета

package genrep

import (
	"time"
)

// Date - тип для работы в json с датами вида 2006-01-02
type Date struct {
	time.Time
}

const dtLayout = "2006-01-02"

// ServiceDef - атрибуты оказываемых услуг
type ServiceDef struct {
	Name  string  `json:"name"`  // наименование
	Units string  `json:"units"` // единица измерения
	Price float64 `json:"price"` // стоимость
}

// Document - атрибуты акта
type Document struct {
	ID          string  `json:"id"`
	Customer    string  `json:"customer"`
	Date        Date    `json:"date"`
	Tax         float64 `json:"tax"`
	ConfirmDate Date    `json:"confirmed"`
	Services    []struct {
		Amount float64 `json:"Amount"`
		Code   string  `json:"Code"`
	} `json:"services"`
}

// Party - сторона акта (заказчик или исполнитель)
type Party struct {
	ID       string `json:"Id"`
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

// CustomerServices - описание услуги для каждого кода
type CustomerServices map[string]ServiceDef

// Act - подготовленные для формирования документа атрибуты акта
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

// ActDef - атрибуты акта для сохранения в json
type ActDef struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Customer    string  `json:"customer"`
	Contractor  string  `json:"contractor"`
	Date        Date    `json:"date"`
	ConfirmDate Date    `json:"confirm"`
	Amount      float64 `json:"amount"`
	Tax         float64 `json:"tax"`
}

// -----------------------------------------------------------------------------

// UnmarshalJSON парсит дату из json
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	d.Time, err = time.Parse(dtLayout, string(b))
	return
}

// MarshalJSON форматирует дату для экспорта в json
func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 12)
	b = append(b, '"')
	b = d.AppendFormat(b, "2006-01-02")
	b = append(b, '"')
	return b, nil
}
