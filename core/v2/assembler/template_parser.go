package assembler

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"

	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
	"github.com/Rhymond/go-money"
)

type TemplateParser struct {
	text     string
	template *template.Template
}

func NewTemplateParser(name string, text string) (*TemplateParser, error) {
	parser := &TemplateParser{text: text}
	tpl, err := template.New(name).Funcs(parser.Funcs()).Parse(text)
	if err != nil {
		return nil, err
	}
	parser.template = tpl
	return parser, nil
}

func (t *TemplateParser) multiplyFloat(a float64, b float64) float64 {
	return a * b
}

func (t *TemplateParser) multiply(a int, b int) int {
	return a * b
}

func (t *TemplateParser) plural(n int, value string) string {
	if n > 1 {
		return strconv.Itoa(n) + value + "s"
	}
	return strconv.Itoa(n) + value
}

func (t *TemplateParser) toFloat(value interface{}) float64 {
	switch val := value.(type) {
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case float64:
		return val
	default:
	}
	return 1.0
}

func (t *TemplateParser) toInt(value interface{}) int {
	switch val := value.(type) {
	case string:
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	case float32:
		return int(val)
	case float64:
		return int(val)
	default:
	}
	return 0
}

func (t *TemplateParser) formatMoney(value interface{}) string {
	var val float64

	switch v := value.(type) {
	case int:
		val = float64(v)
	case float64:
		val = v
	case string:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "-1"
		}
		val = n
	default:
		return ""
	}

	return money.New(int64(val*100.0), "BRL").Display()
}

func (t *TemplateParser) Cleanup(in string) string {
	out := strings.Replace(in, "\r", " ", -1)
	out = strings.Replace(out, `"`, "''", -1)
	return out
}

func (t *TemplateParser) Funcs() template.FuncMap {
	return template.FuncMap{
		"Multiply":            t.multiply,
		"MultiplyFloat":       t.multiplyFloat,
		"FormatMoney":         t.formatMoney,
		"ToInt":               t.toInt,
		"ToFloat":             t.toFloat,
		"Plural":              t.plural,
		"NoMoreRecords":       t.NoMoreRecords,
		"FirstCapitalLetter":  t.FirstCapitalLetter,
		"defaultIntegerValue": t.DefaultIntegerValue,
		"cleanup":             t.Cleanup,
	}
}

func (t *TemplateParser) Parse(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.template.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func cleanTemplate(filledTemplate string) (jsonFinal interface{}, err error) {
	log.Debug("Started", &log.LogContext{"function": "cleanTemplate"})
	defer log.Debug("Finished", &log.LogContext{"function": "cleanTemplate"})

	isCleaned := false

	var oldTrimmedString string

	for !isCleaned {
		oldTrimmedString = filledTemplate
		filledTemplate = strings.ReplaceAll(filledTemplate, "\n", "")
		filledTemplate = strings.ReplaceAll(filledTemplate, "\t", "")
		filledTemplate = strings.ReplaceAll(filledTemplate, ",]", "]")
		filledTemplate = strings.ReplaceAll(filledTemplate, "[,", "[")
		filledTemplate = strings.ReplaceAll(filledTemplate, ",,", ",")

		if oldTrimmedString == filledTemplate {
			isCleaned = true
		}
	}

	err = json.Unmarshal([]byte(filledTemplate), &jsonFinal)
	if err != nil {
		log.Error("Unmarshling template", err, &log.LogContext{
			"template": filledTemplate,
		})
	}

	return
}

func (t *TemplateParser) DefaultIntegerValue(value interface{}) float64 {
	if value == nil {
		return float64(0)
	}

	return value.(float64)
}

func (t *TemplateParser) NoMoreRecords(totalRecords float64, totalRows int) (NoMoreRecords bool) {
	NoMoreRecords = false

	if int(totalRecords) == totalRows {
		NoMoreRecords = true
	}
	return
}

func (t *TemplateParser) FirstCapitalLetter(title string) string {
	res := strings.ToLower(title)
	words := strings.Fields(res)
	smallwords := " de do da "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}
