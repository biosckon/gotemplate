package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// csv or json to map
func file2map(fname string) map[string]string {
	df, err := os.Open(fname)
	if err != nil {
		log.Fatal("Error: opening data file: ", err)
	}
	defer df.Close()

	ret := make(map[string]string)
	ext := filepath.Ext(fname)

	switch ext {
	case ".csv":
		csv_reader := csv.NewReader(df)
		for {
			rec, err := csv_reader.Read()
			if err == io.EOF {
				break
			}
			if len(rec) > 1 {
				ret[rec[0]] = rec[1]
			}

		}

	case ".yaml", ".yml":
		bytes, err := ioutil.ReadAll(df)
		if err != nil {
			log.Fatal("Error reading Yaml file", err)
		}
		err = yaml.Unmarshal(bytes, ret)
		if err != nil {
			log.Fatal("Error parsing Yaml file", err)
		}

	case ".json":
		bytes, err := ioutil.ReadAll(df)
		if err != nil {
			log.Fatal("Error reading json  file", err)
		}
		err = json.Unmarshal(bytes, ret)
		if err != nil {
			log.Fatal("Error parsing json file", err)
		}
	}

	log.Printf("ret map: %v\n", ret)

	return ret
}

func main() {
	flag.Parse()

	// expecting
	// args[0] file with template
	// args[1] file with data file (json, csv); use csv by default
	// if there is a args[2] it is used a output otherwise stdout
	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("Error: Minimum number of arguments 2. Template first followed by the data file (*.json or *.csv).")
	}

	templ_fname := args[0]
	data_fname := args[1]
	out_fname := ""
	outf := os.Stdout

	if len(args) == 3 {
		out_fname = args[2]
	}

	// Parse tempalte file
	tmpl, err := template.ParseFiles(templ_fname)
	if err != nil {
		log.Fatal("Error: parsing template file: ", err)
	}
	tmpl = tmpl.Option("missingkey=error")

	// get data
	data := file2map(data_fname)

	// decide on output file
	if out_fname != "" {
		outf, err = os.Create(out_fname)
		if err != nil {
			log.Fatal("Error: creating/ ovewriting file: ", err)
		}
		defer outf.Close()
	}

	// run the template
	err = tmpl.Execute(outf, data)
	if err != nil {
		log.Fatal("Error excuting template: ", err)
	}
}
