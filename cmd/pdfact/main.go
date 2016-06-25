package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/LeKovr/go-base/logger"
	"github.com/jessevdk/go-flags"

	"github.com/LeKovr/genrep"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Def     string `long:"def" default:"data.json" description:"Document definition json"`
	Key     string `long:"key" default:"default"   description:"Document key in 'document' struct"`
	Out     string `long:"out" default:"act-"      description:".pdf file name prefix"`
	Version bool   `long:"version" description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	log logger.Flags
}

// -----------------------------------------------------------------------------

func main() {

	var cfg Config
	log, _ := setUp(&cfg)
	defer log.Close()

	log.Infof("%s v %s. pdf act generator", path.Base(os.Args[0]), Version)
	log.Println("Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	var def genrep.Act
	file, err := ioutil.ReadFile(cfg.Def)
	if err != nil {
		log.Fatal("Definition read error: ", err)
	}

	err = json.Unmarshal(file, &def)
	if err != nil {
		log.Fatal("Definition parse error: ", err)
	}

	doc, ok := def.Documents[cfg.Key]
	if !ok {
		log.Fatalf("Document %s not found", cfg.Key)
	}

	customer, ok := def.Customers[doc.Customer]
	if !ok {
		log.Fatalf("Customer %s not found", doc.Customer)
	}

	p := cfg.Out + cfg.Key
	if customer.ID != "" {
		prepDir(log, customer.ID)
		p = path.Join(customer.ID, p)
	}
	log.Debugf("Def parsed: %+v", def)
	err = genrep.GenerateAct(def, doc, customer, p)

	if err != nil {
		log.Fatal("Pdf out error: ", err)
	} else {
		fmt.Printf("File %s.pdf generated\n", p)
	}

}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (log *logger.Log, err error) {

	p := flags.NewParser(nil, flags.Default)
	_, err = p.AddGroup("Application Options", "", cfg)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Logging Options", "", &cfg.log)
	panicIfError(err) // check Flags parse error

	_, err = p.Parse()
	if err != nil {
		os.Exit(1) // error message written already
	}

	if cfg.Version {
		// show version & exit
		fmt.Printf("%s\n%s\n%s", Version, Build, Commit)
		os.Exit(0)
	}

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a new instance of the logger
	log, err = logger.New(logger.Dest(cfg.log.LogDest), logger.Level(cfg.log.LogLevel))
	if err != nil {
		panic(err)
	}

	return
}

// -----------------------------------------------------------------------------

func prepDir(log *logger.Log, dir string) {
	if _, err := os.Stat(dir); err == nil {
		// dir exists
		return
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("Dir %s create error: %s", dir, err)
		os.Exit(3)
	}
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
