package custom_attribute

import (
	"fmt"
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/response"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

var columns = []string{"id", "name", "default_value"}

func NewCmdCustomAttribute(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "custom-attribute",
		Aliases: []string{"ca"},
		Short:   "Manage custom attributes",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newSetValueCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all custom attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/custom_attributes"
			params := buildQP(cmd, "limit", "starting_after", "direction")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <attribute-id>",
		Short: "Get custom attribute details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/custom_attributes/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newCreateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom attribute",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			name, _ := cmd.Flags().GetString("name")
			values["name"] = name
			if dv, _ := cmd.Flags().GetString("default-value"); dv != "" {
				values["default_value"] = dv
			}
			data, err := client.DoForm("POST", "/custom_attributes", values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Attribute name (required)")
	cmd.Flags().String("default-value", "", "Default value")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <attribute-id>",
		Short: "Update a custom attribute",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if dv, _ := cmd.Flags().GetString("default-value"); dv != "" {
				values["default_value"] = dv
			}
			data, err := client.DoForm("PATCH", "/custom_attributes/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("default-value", "", "Default value (required)")
	_ = cmd.MarkFlagRequired("default-value")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <attribute-id>",
		Short: "Delete a custom attribute",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_attributes/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Custom attribute deleted.")
			return nil
		},
	}
}

func newSetValueCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-value <attribute-name>",
		Short: "Set attribute value for multiple devices (JSON)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			jsonData, _ := cmd.Flags().GetString("data")
			_, err = client.Do("PUT", "/custom_attribute_values/"+args[0], strings.NewReader(jsonData))
			if err != nil {
				return err
			}
			fmt.Println("Attribute value set.")
			return nil
		},
	}
	cmd.Flags().String("data", "", "JSON array: [{\"device_id\":1,\"value\":\"x\"},...] (required)")
	_ = cmd.MarkFlagRequired("data")
	return cmd
}

func buildQP(cmd *cobra.Command, params ...string) string {
	var parts []string
	for _, p := range params {
		flag := cmd.Flags().Lookup(strings.ReplaceAll(p, "_", "-"))
		if flag == nil {
			continue
		}
		val := flag.Value.String()
		if val != "" && val != "0" {
			parts = append(parts, p+"="+val)
		}
	}
	return strings.Join(parts, "&")
}
