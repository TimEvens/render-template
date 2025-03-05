/*
Copyright (c) 2022 Cisco Systems, Inc. and others.  All rights reserved.
*/
package render_tpl

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"text/template"
)

func Run() {

	loadValues()
	renderTemplate(TemplateFile)

}

func loadValues() {

	for _, f := range ValueFiles {
		log.Printf("File: %s", f)

		Values = coalesceTables(ReadYaml(f, false), Values, f)
	}

	log.Debug(Values)
}

func renderTemplate(filename string) {
	log.Infof("Parsing template: %s", filename)

	tpl := template.Must(template.New(path.Base(filename)).
		Funcs(sprig.TxtFuncMap()).
		Funcs(GetCustomTplFuncMap()).
		ParseFiles(filename))

	var writer *bufio.Writer

	if UseStdout {
		writer = bufio.NewWriter(os.Stdout)

		fmt.Printf("----------[ %s ] ----------------\n", filename)

		defer fmt.Printf("---------- END ---------------------------------\n")

	} else {
		outFilename := strings.Replace(filename, ".tpl", "", 1)

		output, err := os.Create(outFilename)
		if err != nil {
			log.Fatal("Unable to create output file ", outFilename)
		}

		writer = bufio.NewWriter(output)

		defer output.Close()
	}

	err := tpl.Execute(writer, Values)
	if err != nil {
		log.Fatal("Error executing ", err)
	}

	defer writer.Flush()
}
