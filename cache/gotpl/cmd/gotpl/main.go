/*
Copyright (c) 2023 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/pflag"
)

var (
	tplfile string
	varfile string
)

func main() {
	errlog := log.New(os.Stderr, "", 0)

	pflag.Usage = func() {
		errlog.Printf("Usage: %s [options] [template]\n", os.Args[0])
		errlog.Printf("[template]: path to go template to be rendered; optional; if missing or '-', template will be read from stdin\n")
		errlog.Printf("[options]:\n")
		pflag.PrintDefaults()
	}
	pflag.StringVarP(&varfile, "var_file", "f", "", "Path to a file containing binding variables in JSON or YAML format")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	tplfile = pflag.Arg(0)

	var rawtpl []byte
	if tplfile == "" || tplfile == "-" {
		raw, err := readStdin()
		if err != nil {
			errlog.Fatal(err)
		}
		rawtpl = raw
	} else {
		raw, err := readFile(tplfile)
		if err != nil {
			errlog.Fatal(err)
		}
		rawtpl = raw
	}

	var rawvar []byte
	if varfile == "" {
		errlog.Fatal("flag --var_file empty or not provided")
	} else {
		raw, err := readFile(varfile)
		if err != nil {
			errlog.Fatal(err)
		}
		rawvar = raw
	}

	tpl := string(rawtpl)

	var data map[string]interface{}
	if err := json.Unmarshal(rawvar, &data); err != nil {
		errlog.Fatal(err)
	}

	tmpl, err := template.New("main").Funcs(sprig.TxtFuncMap()).Option("missingkey=error").Parse(tpl)
	if err != nil {
		errlog.Fatal(err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		errlog.Fatal(err)
	}

	fmt.Println(out.String())
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func readStdin() ([]byte, error) {
	raw, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
