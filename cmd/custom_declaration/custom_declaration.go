package custom_declaration

import (
	"fmt"
	"os"
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/response"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

var columns = []string{"id", "name", "declaration_type"}

func NewCmdCustomDeclaration(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "custom-declaration",
		Aliases: []string{"cd"},
		Short:   "Manage custom declarations",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newDownloadCmd(opts))
	cmd.AddCommand(newPushToDeviceCmd(opts))
	cmd.AddCommand(newRemoveFromDeviceCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all custom declarations",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/custom_declarations"
			params := buildQP(cmd, "limit", "starting_after", "direction", "search")
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
	cmd.Flags().String("search", "", "Search by name")
	return cmd
}

func newCreateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom declaration",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			name, _ := cmd.Flags().GetString("name")
			values["name"] = name
			dt, _ := cmd.Flags().GetString("declaration-type")
			values["declaration_type"] = dt
			if us, _ := cmd.Flags().GetString("user-scope"); us != "" {
				values["user_scope"] = us
			}
			if cmd.Flags().Changed("attribute-support") {
				as, _ := cmd.Flags().GetBool("attribute-support")
				if as {
					values["attribute_support"] = "true"
				} else {
					values["attribute_support"] = "false"
				}
			}
			if cmd.Flags().Changed("escape-attributes") {
				ea, _ := cmd.Flags().GetBool("escape-attributes")
				if ea {
					values["escape_attributes"] = "true"
				} else {
					values["escape_attributes"] = "false"
				}
			}
			if ap, _ := cmd.Flags().GetString("activation-predicate"); ap != "" {
				values["activation_predicate"] = ap
			}
			filePath, _ := cmd.Flags().GetString("payload")
			data, err := client.DoMultipart("POST", "/custom_declarations", values, "payload", filePath)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Declaration name (required)")
	cmd.Flags().String("declaration-type", "", "Declaration type (required)")
	cmd.Flags().String("payload", "", "Path to payload file (required)")
	cmd.Flags().String("user-scope", "", "User scope")
	cmd.Flags().Bool("attribute-support", false, "Enable attribute support")
	cmd.Flags().Bool("escape-attributes", false, "Escape attributes")
	cmd.Flags().String("activation-predicate", "", "Activation predicate")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("declaration-type")
	_ = cmd.MarkFlagRequired("payload")
	return cmd
}

func newDownloadCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <declaration-id>",
		Short: "Download a custom declaration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			body, filename, err := client.DoDownload("/custom_declarations/" + args[0] + "/download")
			if err != nil {
				return err
			}
			outPath, _ := cmd.Flags().GetString("output")
			if outPath == "" {
				if filename != "" {
					outPath = filename
				} else {
					outPath = "declaration_" + args[0] + ".json"
				}
			}
			if err := os.WriteFile(outPath, body, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("Declaration downloaded to %s\n", outPath)
			return nil
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output file path")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <declaration-id>",
		Short: "Update a custom declaration",
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
			if dt, _ := cmd.Flags().GetString("declaration-type"); dt != "" {
				values["declaration_type"] = dt
			}
			if us, _ := cmd.Flags().GetString("user-scope"); us != "" {
				values["user_scope"] = us
			}
			if cmd.Flags().Changed("attribute-support") {
				as, _ := cmd.Flags().GetBool("attribute-support")
				if as {
					values["attribute_support"] = "true"
				} else {
					values["attribute_support"] = "false"
				}
			}
			if cmd.Flags().Changed("escape-attributes") {
				ea, _ := cmd.Flags().GetBool("escape-attributes")
				if ea {
					values["escape_attributes"] = "true"
				} else {
					values["escape_attributes"] = "false"
				}
			}
			if ap, _ := cmd.Flags().GetString("activation-predicate"); ap != "" {
				values["activation_predicate"] = ap
			}
			filePath, _ := cmd.Flags().GetString("payload")
			if filePath != "" {
				data, err := client.DoMultipart("PATCH", "/custom_declarations/"+args[0], values, "payload", filePath)
				if err != nil {
					return err
				}
				return response.Print(opts.GetOutputFormat(), data, columns)
			}
			data, err := client.DoForm("PATCH", "/custom_declarations/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Declaration name")
	cmd.Flags().String("declaration-type", "", "Declaration type")
	cmd.Flags().String("payload", "", "Path to payload file")
	cmd.Flags().String("user-scope", "", "User scope")
	cmd.Flags().Bool("attribute-support", false, "Enable attribute support")
	cmd.Flags().Bool("escape-attributes", false, "Escape attributes")
	cmd.Flags().String("activation-predicate", "", "Activation predicate")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <declaration-id>",
		Short: "Delete a custom declaration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_declarations/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Custom declaration deleted.")
			return nil
		},
	}
}

func newPushToDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "push-to-device <declaration-id> <device-id>",
		Short: "Push declaration to device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/custom_declarations/"+args[0]+"/devices/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Declaration pushed to device.")
			return nil
		},
	}
}

func newRemoveFromDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-from-device <declaration-id> <device-id>",
		Short: "Remove declaration from device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_declarations/" + args[0] + "/devices/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("Declaration removed from device.")
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
