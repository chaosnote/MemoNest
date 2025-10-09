package utils

import (
	"html/template"
)

type TemplateConfig struct {
	Layout  string   // 空字串則不使用
	Page    []string // []string{name.html,...}
	Pattern []string // []string{*.html}
}

func RenderTemplate(config TemplateConfig) (tmpl *template.Template, e error) {
	list := []string{}
	if len(config.Layout) > 0 {
		list = append(list, config.Layout)
	}
	list = append(list, config.Page...)

	tmpl, e = template.ParseFiles(list...)
	if e != nil {
		return
	}
	for _, value := range config.Pattern {
		tmpl, e = tmpl.ParseGlob(value)
		if e != nil {
			return
		}
	}
	return
}
