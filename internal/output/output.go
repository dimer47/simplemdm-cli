package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

func Print(format string, data interface{}, columns []string) error {
	switch format {
	case "json":
		return printJSON(data)
	case "yaml":
		return printYAML(data)
	case "csv":
		return printCSV(data, columns)
	default:
		return printTable(data, columns)
	}
}

func printJSON(data interface{}) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func printYAML(data interface{}) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}

func printTable(data interface{}, columns []string) error {
	rows := normalize(data)
	if len(rows) == 0 {
		fmt.Println("No results.")
		return nil
	}

	if len(columns) == 0 {
		columns = inferColumns(rows)
	}

	table := tablewriter.NewTable(os.Stdout)

	headers := make([]any, len(columns))
	for i, col := range columns {
		headers[i] = strings.ToUpper(col)
	}
	table.Header(headers...)

	for _, row := range rows {
		vals := make([]any, len(columns))
		for i, col := range columns {
			if v, ok := row[col]; ok {
				vals[i] = fmt.Sprintf("%v", v)
			}
		}
		table.Append(vals...)
	}

	table.Render()
	return nil
}

func printCSV(data interface{}, columns []string) error {
	rows := normalize(data)
	if len(rows) == 0 {
		return nil
	}

	if len(columns) == 0 {
		columns = inferColumns(rows)
	}

	w := csv.NewWriter(os.Stdout)
	_ = w.Write(columns)

	for _, row := range rows {
		vals := make([]string, len(columns))
		for i, col := range columns {
			if v, ok := row[col]; ok {
				vals[i] = fmt.Sprintf("%v", v)
			}
		}
		_ = w.Write(vals)
	}
	w.Flush()
	return nil
}

func normalize(data interface{}) []map[string]interface{} {
	switch v := data.(type) {
	case []map[string]interface{}:
		return v
	case map[string]interface{}:
		return []map[string]interface{}{v}
	case []interface{}:
		var result []map[string]interface{}
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	default:
		if data == nil {
			return nil
		}
		b, err := json.Marshal(data)
		if err != nil {
			return nil
		}
		var result interface{}
		if err := json.Unmarshal(b, &result); err != nil {
			return nil
		}
		// Avoid infinite recursion if round-trip doesn't produce a known type
		switch result.(type) {
		case map[string]interface{}, []interface{}, []map[string]interface{}:
			return normalize(result)
		default:
			return nil
		}
	}
}

func inferColumns(rows []map[string]interface{}) []string {
	seen := make(map[string]bool)
	var cols []string
	for _, row := range rows {
		for k := range row {
			if !seen[k] {
				seen[k] = true
				cols = append(cols, k)
			}
		}
	}
	sort.Strings(cols)
	return cols
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
}

func MaskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return "****" + token[len(token)-8:]
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func FormatList(items []string) string {
	return strings.Join(items, ", ")
}
