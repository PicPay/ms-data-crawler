package configuration

import (
	"bytes"
	"text/template"
)

type TemplateParser struct {
	text     string
	template *template.Template
}

func NewTemplateParser(name string, text string) (*TemplateParser, error) {
	parser := &TemplateParser{text: text}
	tpl, err := template.New(name).Funcs(parser.getFunctionList()).Parse(text)
	if err != nil {
		return nil, err
	}
	parser.template = tpl
	return parser, nil
}

func (t *TemplateParser) Parse(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.template.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TemplateParser) getFunctionList() template.FuncMap {
	return template.FuncMap{
		"HasMoreItems": t.HasMoreItems,
	}
}

func (t *TemplateParser) HasMoreItems(elementPosition int, totalRows int) (NoMoreItems bool) {
	NoMoreItems = true

	if (elementPosition + 1) == totalRows {
		NoMoreItems = false
	}
	return
}
