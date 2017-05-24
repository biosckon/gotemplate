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

// csv/ json/ yaml to map
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
		err = yaml.Unmarshal(bytes, &ret)
		if err != nil {
			log.Fatal("Error parsing Yaml file", err)
		}

	case ".json":
		bytes, err := ioutil.ReadAll(df)
		if err != nil {
			log.Fatal("Error reading json  file", err)
		}
		err = json.Unmarshal(bytes, &ret)
		if err != nil {
			log.Fatal("Error parsing json file", err)
		}
	}

	return ret
}

func init_tpl(globs []string) *template.Template {
	// collect all files, extend globs
	files := make(map[string]bool)

	for _, glob := range globs {
		fs, err := filepath.Glob(glob)
		if err != nil {
			log.Println("Error parsing file as Glob: ", err)
			continue
		}

		for _, f := range fs {
			files[f] = true
		}
	}

	all_files := make([]string, len(files))

	if len(all_files) == 0 {
		log.Fatal("init_tpl has no files to process")
	}

	i := 0
	for f := range files {
		all_files[i] = f
		i++
	}

	tpl := template.Must(template.ParseFiles(all_files...))
	tpl = tpl.Option("missingkey=error")

	return tpl
}

func main() {
	flag.Parse()

	// expecting at least 2 arguments
	// last argument is name of the data file
	// data file in (yaml/yml, json, csv) format
	// 0 to all but last are expected to be templates
	// all output shall go to stdout
	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("Error: Minimum number of arguments 2. Template first followed by the data file (*.json or *.csv) last.")
	}

	_last := len(args) - 1
	tpl_fnames := args[0:_last] // all but last
	data_fname := args[_last]
	outf := os.Stdout
	var err error

	// Parse tempalte files
	tpl := init_tpl(tpl_fnames)

	// get data
	data := file2map(data_fname)

	// run the template
	err = tpl.ExecuteTemplate(outf, "index", data)
	if err != nil {
		log.Fatal("Error excuting template: ", err)
	}
}
