package genrep

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"reflect"

	"encoding/json"
	"io/ioutil"

	"github.com/LeKovr/num2word"
)

// Convert 'ABCDEFG' to, for example, 'A,BCD,EFG'
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

func getPartyField(p *Party, field string) string {
	r := reflect.ValueOf(p)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

var monthName = []string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}

func dateName(date Date) string {
	d := date.Time
	return fmt.Sprintf("%d %s %d", d.Day(), monthName[d.Month()-1], d.Year())
}

var fields = []string{"Address", "Rs", "Ks", "Bank", "Bik", "Phone"}

var label = map[string]string{
	"Inn":     "ИНН",
	"Kpp":     "КПП",
	"Address": "Адрес",
	"Rs":      "Р/с",
	"Ks":      "К/с",
	"Bank":    "Банк",
	"Bik":     "БИК",
	"Phone":   "Телефон",
}

var header = []string{"№", "Наименование работ, услуг", "Кол-во", "Ед.", "Цена", "Сумма"}

// Column widths
var w = []float64{7.0, 97.0, 15.0, 15.0, 20.0, 22.0}

// GenerateAct generates out file with act
func GenerateAct(def Act, doc Document, customer Party, out string) (err error) {

	pdf := gofpdf.New("P", "mm", "A4", "fonts")
	pdf.SetCellMargin(2)
	pdf.SetFillColor(94, 94, 94)
	pdf.SetDrawColor(145, 145, 145)
	pdf.AddPage()

	pdf.AddFont("Helvetica", "", def.Meta.Font+".json")
	pdf.AddFont("Helvetica", "B", def.Meta.Font+"_Bold.json")

	fontSize := 11.0
	pdf.SetFont("Helvetica", "B", fontSize)
	_, lineHt := pdf.GetFontSize()

	pdf.SetLeftMargin(15)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	title := fmt.Sprintf(def.Names.TitleFmt, doc.ID, dateName(doc.Date))
	pdf.CellFormat(175, lineHt, tr(title), "", 0, "C", false, 0, "")
	pdf.Ln(lineHt)

	fontSize = 10.0
	pdf.SetFont("Helvetica", "", fontSize)

	pdf.CellFormat(175, lineHt, "", "B", 0, "", false, 0, "")
	pdf.Ln(lineHt * 2)

	pdf.CellFormat(27, lineHt, tr(def.Contractor.Title), "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "B", fontSize)
	pdf.CellFormat(149, lineHt, tr(def.Contractor.NameLong), "", 0, "", false, 0, "")
	pdf.Ln(lineHt + 2)

	pdf.SetFont("Helvetica", "", fontSize)
	pdf.CellFormat(27, lineHt, tr(customer.Title), "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "B", fontSize)
	pdf.CellFormat(149, lineHt, tr(customer.NameLong), "", 0, "", false, 0, "")
	pdf.Ln(lineHt * 3)

	pdf.SetFont("Helvetica", "", fontSize)

	// 	Header
	for j, str := range header {
		pdf.CellFormat(w[j], lineHt+2, tr(str), "1", 0, "CM", false, 0, "")
	}
	pdf.Ln(-1)

	services := def.Services[doc.Customer]
	total := 0.0
	_, lineHt = pdf.GetFontSize()
	for i, srv := range doc.Services {

		srvDef := services[srv.Code]

		lines := pdf.SplitLines([]byte(tr(srvDef.Name)), w[1])
		ht := float64(len(lines)) * (lineHt + 2)

		pdf.CellFormat(w[0], ht, fmt.Sprintf("%d", i+1), "1", 0, "CM", false, 0, "")

		x, y := pdf.GetXY()
		pdf.MultiCell(w[1], lineHt+2, tr(srvDef.Name), "1", "", false)
		pdf.MoveTo(x+w[1], y)

		pdf.CellFormat(w[2], ht, fmt.Sprintf("%.0f", srv.Amount), "1", 0, "RM", false, 0, "")
		pdf.CellFormat(w[3], ht, tr(srvDef.Units), "1", 0, "CM", false, 0, "")
		pdf.CellFormat(w[4], ht, strDelimit(fmt.Sprintf("%.f", srvDef.Price), " ", 3)+".00", "1", 0, "RM", false, 0, "")
		pdf.CellFormat(w[5], ht, strDelimit(fmt.Sprintf("%.0f", srv.Amount*srvDef.Price), " ", 3)+".00", "1", 0, "RM", false, 0, "")
		total += srv.Amount * srvDef.Price
		pdf.Ln(ht)

	}
	pdf.Ln(lineHt * 2)
	pdf.SetFont("Helvetica", "B", fontSize)
	pdf.CellFormat(154, lineHt, tr(def.Names.TotalTitle), "", 0, "RM", false, 0, "")
	pdf.CellFormat(22, lineHt, strDelimit(fmt.Sprintf("%.0f", total), " ", 3)+".00", "", 0, "RM", false, 0, "")
	pdf.Ln(lineHt + 2)

	pdf.SetFont("Helvetica", "", fontSize)
	pdf.CellFormat(154, lineHt, tr(def.Names.TotalTax), "", 0, "RM", false, 0, "")
	pdf.CellFormat(22, lineHt, tr(def.Names.ZeroTax), "", 0, "RM", false, 0, "")

	pdf.Ln(lineHt * 3)

	pdf.MultiCell(175, lineHt+2, tr(fmt.Sprintf(def.Names.TotalWordsFmt, num2word.RuMoney(total, true))), "", "", false)

	pdf.Ln(3)

	pdf.CellFormat(175, 0.6, "", "B", 0, "", true, 0, "")
	pdf.Ln(lineHt * 2)

	fontSize = 9.0
	pdf.SetFont("Helvetica", "", fontSize)
	_, lineHt = pdf.GetFontSize()

	pdf.CellFormat(160, lineHt, tr(def.Names.Resume), "", 0, "", false, 0, "")
	pdf.Ln(lineHt * 2)

	fontSize = 10
	pdf.SetFont("Helvetica", "", fontSize)
	_, lineHt = pdf.GetFontSize()

	lWidth := 23.0
	baseY := pdf.GetY()
	for idx, c := range []Party{def.Contractor, customer} {
		pdf.CellFormat(lWidth+1, lineHt, tr(c.Title), "", 0, "T", false, 0, "")
		pdf.SetFont("Helvetica", "B", fontSize)
		pdf.MultiCell(63, lineHt, tr(c.NameFull), "", "T", false)
		pdf.Ln(-1)

		pdf.SetY(baseY + 10.0)

		pdf.SetFont("Helvetica", "", fontSize)

		pdf.CellFormat(lWidth, lineHt, tr(label["Inn"]), "", 0, "", false, 0, "")
		pdf.CellFormat(27, lineHt, tr(c.Inn), "B", 0, "", false, 0, "")

		pdf.CellFormat(12, lineHt, tr(label["Kpp"]), "", 0, "", false, 0, "")
		pdf.CellFormat(23, lineHt, tr(c.Kpp), "B", 0, "", false, 0, "")
		pdf.Ln(lineHt + 3)

		for _, n := range fields {
			y := pdf.GetY()
			pdf.CellFormat(lWidth, lineHt, tr(label[n]), "", 0, "T", false, 0, "")
			pdf.MultiCell(60, lineHt+1, tr(getPartyField(&c, n)), "B", "T", false)
			pdf.Ln(-1)
			if n == "Address" || n == "Bank" {
				pdf.SetY(y + lineHt*5)
			}
		}

		pdf.CellFormat(25, lineHt*2, "", "B", 0, "", false, 0, "")
		pdf.CellFormat(8, lineHt*2, "/", "", 0, "C", false, 0, "")

		lines := pdf.SplitLines([]byte(tr(c.Name)), 50)
		ht := float64(3-len(lines)) * (lineHt)
		pdf.MultiCell(50, ht, tr(c.Name), "B", "CB", false)

		pdf.Ln(lineHt - 2)
		pdf.CellFormat(25, lineHt, tr(def.Names.Stamp), "", 0, "C", false, 0, "")
		pdf.CellFormat(8, lineHt, "", "", 0, "C", false, 0, "")

		pdf.SetFont("Helvetica", "", fontSize-2)
		pdf.CellFormat(50, lineHt, tr(def.Names.Sign), "", 0, "C", false, 0, "")
		pdf.SetFont("Helvetica", "", fontSize)

		m := 105.0 * float64(idx+1)
		pdf.SetLeftMargin(m)
		pdf.SetXY(m, baseY)

	}
	dt := doc.Date.Format("2006-01-02")
	pdf.SetTitle(fmt.Sprintf(def.Meta.TitleFmt, customer.NameFull, doc.ID, dateName(doc.Date)), true)
	pdf.SetAuthor(fmt.Sprintf(def.Meta.AuthorFmt, def.Contractor.NameFull), true)
	pdf.SetSubject(fmt.Sprintf(def.Meta.SubjectFmt, customer.NameFull, doc.ID, dateName(doc.Date), total), true)
	pdf.SetKeywords(fmt.Sprintf("%s %s %s %s %s %.2f", def.Meta.Keywords, def.Contractor.NameFull, customer.NameFull, doc.ID, dt, total), true)
	pdf.SetCreator(def.Meta.Creator, true)

	err = pdf.OutputFileAndClose(out + ".pdf")

	if err == nil {
		// save .json act definition if no errors
		a := ActDef{
			ID:          doc.ID,
			Title:       title,
			Customer:    customer.NameFull,
			Contractor:  def.Contractor.NameFull,
			Date:        doc.Date,
			Amount:      total,
			Tax:         doc.Tax,
			ConfirmDate: doc.ConfirmDate,
		}
		outjs, _ := json.MarshalIndent(a, "   ", "   ")
		err = ioutil.WriteFile(out+".json", outjs, 0644)
	}

	return
}
