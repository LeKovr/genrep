package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"it.elfire.ru/elfire/reports"

)

var cfgDef, cfgKey, cfgOut string

func init() {
	flag.StringVar(&cfgDef, "def", "data.json", "Document definition json")
	flag.StringVar(&cfgKey, "key", "default", "Document key in 'document' struct")
	flag.StringVar(&cfgOut, "out", "act-", ".pdf file name prefix")
}

func main() {

	log.Println("elact v 1.1, Elfire Act generator")
	log.Println("Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "\nThis is a Act generator for json data")
		fmt.Fprintf(os.Stderr, "\nUsage:\n  elact [options]\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	var def reports.Act
	file, err := ioutil.ReadFile(cfgDef)
	if err != nil {
		log.Println("Definition read error:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &def)
	if err != nil {
		log.Println("Definition parse error:", err)
		os.Exit(2)
	}

	doc, ok := def.Documents[cfgKey]
	if !ok {
		log.Printf("Document %s not found", cfgKey)
		os.Exit(3)
	}

	customer, ok := def.Customers[doc.Customer]
	if !ok {
		log.Printf("Customer %s not found", doc.Customer)
		os.Exit(3)
	}

	//log.Printf("Def parsed: %+v", def)
	err = reports.GenerateAct(def, doc, customer, cfgOut+cfgKey)

	if err != nil {
		log.Println("Pdf out error:", err)
	} else {
		log.Printf("Act %s%s.pdf generated", cfgOut, cfgKey)
	}

}
