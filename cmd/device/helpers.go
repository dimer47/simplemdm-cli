package device

import (
	"encoding/json"
	"fmt"

	"github.com/dimer47/simplemdm-cli/internal/output"
)

func printResponse(format string, data []byte, columns []string) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println(string(data))
		return nil
	}
	if d, ok := raw["data"]; ok {
		if attrs, ok := d.(map[string]interface{}); ok {
			row := flattenItem(attrs)
			return output.Print(format, row, columns)
		}
		if arr, ok := d.([]interface{}); ok {
			var rows []map[string]interface{}
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					rows = append(rows, flattenItem(m))
				}
			}
			return output.Print(format, rows, columns)
		}
	}
	return output.Print(format, raw, columns)
}

func printListResponse(format string, data []byte, columns []string) error {
	return printResponse(format, data, columns)
}

func flattenItem(m map[string]interface{}) map[string]interface{} {
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
	return result
}
