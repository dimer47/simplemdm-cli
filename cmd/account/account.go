package account

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/output"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

func NewCmdAccount(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage SimpleMDM account",
	}
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get account details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/account")
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, []string{"id", "name", "apple_store_country_code"})
		},
	}
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update account details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				values["name"] = name
			}
			if country, _ := cmd.Flags().GetString("apple-store-country-code"); country != "" {
				values["apple_store_country_code"] = country
			}
			if len(values) == 0 {
				return fmt.Errorf("at least one flag (--name or --apple-store-country-code) is required")
			}
			data, err := client.DoForm("PATCH", "/account", values)
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, []string{"id", "name", "apple_store_country_code"})
		},
	}
	cmd.Flags().String("name", "", "Account name")
	cmd.Flags().String("apple-store-country-code", "", "Apple Store country code")
	return cmd
}

func printResponse(format string, data []byte, columns []string) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		// Try printing raw
		fmt.Println(strings.TrimSpace(string(data)))
		return nil
	}
	if d, ok := raw["data"]; ok {
		if attrs, ok := d.(map[string]interface{}); ok {
			if a, ok := attrs["attributes"]; ok {
				if attrMap, ok := a.(map[string]interface{}); ok {
					if id, ok := attrs["id"]; ok {
						attrMap["id"] = id
					}
					return output.Print(format, attrMap, columns)
				}
			}
			return output.Print(format, attrs, columns)
		}
		if arr, ok := d.([]interface{}); ok {
			var rows []map[string]interface{}
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					row := flattenItem(m)
					rows = append(rows, row)
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
