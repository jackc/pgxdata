package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

func initCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "init requires exactly one argument")
		os.Exit(1)
	}

	data := initData{
		PkgName: args[0],
		Version: VERSION,
	}

	err := os.Mkdir(data.PkgName, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	templates := loadTemplates()

	files := []struct {
		path string
		tmpl *template.Template
	}{
		{"config.toml", templates.Lookup("config")},
		{"pgxdata_db.go", templates.Lookup("db")},
	}
	for _, f := range files {
		err := writeInitFile(filepath.Join(data.PkgName, f.path), f.tmpl, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func writeInitFile(path string, tmpl *template.Template, data initData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
