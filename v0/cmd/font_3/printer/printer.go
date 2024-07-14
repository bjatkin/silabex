package printer

import (
	"html/template"
	"strings"
)

type Printer struct {
}

type html struct {
	Runes []string
}

func Render(s string) (string, error) {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		return "", err
	}

	ret := &strings.Builder{}
	err = tmpl.Execute(ret, html{})
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}
