package script

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

var columns = []string{"id", "name", "variable_support", "created_at"}

func NewCmdScript(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "Manage scripts",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/scripts"
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
		Use:   "get <script-id>",
		Short: "Get script details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/scripts/" + args[0])
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
		Short: "Create a script",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			name, _ := cmd.Flags().GetString("name")
			values["name"] = name
			if cmd.Flags().Changed("variable-support") {
				vs, _ := cmd.Flags().GetBool("variable-support")
				if vs {
					values["variable_support"] = "1"
				} else {
					values["variable_support"] = "0"
				}
			}
			filePath, _ := cmd.Flags().GetString("file")
			data, err := client.DoMultipart("POST", "/scripts", values, "file", filePath)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Script name (required)")
	cmd.Flags().String("file", "", "Path to script file (required)")
	cmd.Flags().Bool("variable-support", false, "Enable variable support")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <script-id>",
		Short: "Update a script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				values["name"] = name
			}
			if cmd.Flags().Changed("variable-support") {
				vs, _ := cmd.Flags().GetBool("variable-support")
				if vs {
					values["variable_support"] = "true"
				} else {
					values["variable_support"] = "false"
				}
			}
			filePath, _ := cmd.Flags().GetString("file")
			if filePath != "" {
				data, err := client.DoMultipart("PATCH", "/scripts/"+args[0], values, "file", filePath)
				if err != nil {
					return err
				}
				return response.Print(opts.GetOutputFormat(), data, columns)
			}
			data, err := client.DoForm("PATCH", "/scripts/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Script name")
	cmd.Flags().String("file", "", "Path to script file")
	cmd.Flags().Bool("variable-support", false, "Enable variable support")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <script-id>",
		Short: "Delete a script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/scripts/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Script deleted.")
			return nil
		},
	}
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
