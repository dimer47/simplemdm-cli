package response

import (
	"encoding/json"
	"fmt"

	"github.com/dimer47/simplemdm-cli/internal/output"
)

func Print(format string, data []byte, columns []string) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println(string(data))
		return nil
	}
	if d, ok := raw["data"]; ok {
		switch v := d.(type) {
		case map[string]interface{}:
			return output.Print(format, FlattenItem(v), columns)
		case []interface{}:
			var rows []map[string]interface{}
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					rows = append(rows, FlattenItem(m))
				}
			}
			return output.Print(format, rows, columns)
		}
	}
	return output.Print(format, raw, columns)
}

func FlattenItem(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if id, ok := m["id"]; ok {
		result["id"] = id
	}
	if t, ok := m["type"]; ok {
		result["type"] = t
	}
	if attrs, ok := m["attributes"]; ok {
		if attrMap, ok := attrs.(map[string]interface{}); ok {
			for k, v := range attrMap {
				result[k] = v
			}
		}
	}
	// If no attributes key, copy all top-level fields
	if _, hasAttrs := m["attributes"]; !hasAttrs {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func PrintMessage(msg string) {
	fmt.Println(msg)
}
