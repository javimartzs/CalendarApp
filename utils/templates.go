package utils

import (
	"html/template"
	"log"
	"strings"
	"time"
)

var Tmpl *template.Template

// Funcion para iniciar las plantillas hmtl
func InitTemplates() {
	funcMap := template.FuncMap{
		"seq": seq,
		"add": add,
	}

	var err error
	Tmpl, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

// Funcion para convertir los meses ingleses a castellano
func TranslateMonth(date string) string {
	months := map[string]string{
		"January":   "Enero",
		"February":  "Febrero",
		"March":     "Marzo",
		"April":     "Abril",
		"May":       "Mayo",
		"June":      "Junio",
		"July":      "Julio",
		"August":    "Agosto",
		"September": "Septiembre",
		"October":   "Octubre",
		"November":  "Noviembre",
		"December":  "Diciembre",
	}

	for en, es := range months {
		date = strings.Replace(date, en, es, 1)
	}
	return date
}

// Funcion para transformar fechas a string
func FormatDate(t time.Time) string {
	return TranslateMonth(t.Format("2 de January"))
}

// seq generates a sequence of integers from start to end inclusive
func seq(start, end int) []int {
	var s []int
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}

// add returns the sum of two integers
func add(a, b int) int {
	return a + b
}
