package utils

import (
	"html/template"
)

type TemplateConfig struct {
	Layout  string           // 空字串則不使用
	Page    []string         // []string{name.html,...}
	Pattern []string         // []string{*.html}
	Funcs   template.FuncMap // 外部函式註入
}

func RenderTemplate(config TemplateConfig) (tmpl *template.Template, e error) {
	list := []string{}
	if len(config.Layout) > 0 {
		list = append(list, config.Layout)
	}
	list = append(list, config.Page...)

	// 建立 base template 並注入 FuncMap（如果有）
	if config.Funcs != nil {
		tmpl = template.New("base").Funcs(config.Funcs)
	} else {
		tmpl = template.New("base")
	}

	tmpl, e = tmpl.ParseFiles(list...)
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
