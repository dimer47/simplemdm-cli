package device_group

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

var columns = []string{"id", "name"}

func NewCmdDeviceGroup(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "device-group",
		Aliases: []string{"dg"},
		Short:   "Manage device groups",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newAssignDeviceCmd(opts))
	cmd.AddCommand(newCloneCmd(opts))
	cmd.AddCommand(newCustomAttributesCmd(opts))
	cmd.AddCommand(newSetCustomAttributeCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all device groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/device_groups"
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
		Use:   "get <group-id>",
		Short: "Get device group details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/device_groups/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newAssignDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "assign-device <group-id> <device-id>",
		Short: "Assign device to group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/device_groups/"+args[0]+"/devices/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Device assigned to group.")
			return nil
		},
	}
}

func newCloneCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clone <group-id>",
		Short: "Clone a device group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Post("/device_groups/"+args[0]+"/clone", nil)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newCustomAttributesCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "custom-attributes <group-id>",
		Short: "Get custom attribute values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/device_groups/" + args[0] + "/custom_attribute_values")
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newSetCustomAttributeCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-custom-attribute <group-id> <attribute-name>",
		Short: "Set a custom attribute value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			value, _ := cmd.Flags().GetString("value")
			_, err = client.DoForm("PUT", "/device_groups/"+args[0]+"/custom_attribute_values/"+args[1], map[string]string{"value": value})
			if err != nil {
				return err
			}
			fmt.Println("Custom attribute set.")
			return nil
		},
	}
	cmd.Flags().String("value", "", "Attribute value (required)")
	_ = cmd.MarkFlagRequired("value")
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
