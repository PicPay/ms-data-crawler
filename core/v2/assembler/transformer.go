// Trasformer handles a response body transformations to one or more functions

package assembler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

type TransformerFn map[string]func(body []byte, args []string) []byte

var ResponseTransformer TransformerFn = TransformerFn{
	// returns the first item from array
	"first": func(body []byte, args []string) []byte {
		data, _, _, err := jsonparser.Get(body, "[0]")
		if err != nil {
			return body
		}
		return data
	},

	// filter removes fields from json, leaving only those in arguments
	"filter": func(body []byte, args []string) []byte {
		var b map[string]interface{}

		err := json.Unmarshal(body, &b)
		if err != nil {
			return body
		}

		data := make(map[string]interface{})
		for _, arg := range args {
			val, ok := b[arg]
			if ok {
				data[arg] = val
			}
		}

		final, err := json.Marshal(data)
		if err != nil {
			return body
		}

		return final
	},

	// iterate over N times and replicate the object, appends an `i` for index starting from 1
	"iterate": func(body []byte, args []string) []byte {
		if len(args) < 1 {
			return body
		}

		n, err := strconv.Atoi(args[0])
		if err != nil {
			return body
		}

		var data []map[string]interface{}
		for i := 0; i < n; i++ {
			var b map[string]interface{}

			err := json.Unmarshal(body, &b)
			if err != nil {
				return body
			}

			b["i"] = i + 1
			data = append(data, b)
		}

		final, err := json.Marshal(data)
		if err != nil {
			return body
		}

		return final
	},

	// group array by some field
	"group_by": func(body []byte, args []string) []byte {
		var bodyParsed []map[string]interface{}
		if err := json.Unmarshal(body, &bodyParsed); err != nil {
			return body
		}

		if len(args) < 1 {
			return body
		}

		key := args[0]
		final := make(map[string][]map[string]interface{}, 0)
		for _, item := range bodyParsed {
			value, found := item[key]
			if !found {
				continue
			}

			// the value of the field matching "key" var
			var fieldValue string
			switch value.(type) {
			case string:
				fieldValue = value.(string)
			case []string:
				values := value.([]string)
				if len(values) > 0 {
					fieldValue = values[0]
				}
			case []interface{}:
				values := value.([]interface{})
				if len(values) > 0 {
					val, ok := values[0].(string)
					if ok {
						fieldValue = val
					}
				}
			default:
			}

			fieldValue = strings.ReplaceAll(fieldValue, "/", "")
			fieldValue = strings.ReplaceAll(fieldValue, "\\", "")
			if fieldValue == "" {
				continue
			}

			if _, ok := final[fieldValue]; !ok {
				final[fieldValue] = make([]map[string]interface{}, 0)
			}
			final[fieldValue] = append(final[fieldValue], item)
		}

		if data, err := json.Marshal(final); err == nil {
			return data
		}

		return body
	},
}
