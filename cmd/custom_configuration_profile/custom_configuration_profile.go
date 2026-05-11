package custom_configuration_profile

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

var columns = []string{"id", "name", "user_scope"}

func NewCmdCustomConfigurationProfile(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "custom-configuration-profile",
		Aliases: []string{"ccp"},
		Short:   "Manage custom configuration profiles",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newDownloadCmd(opts))
	cmd.AddCommand(newPushToDeviceCmd(opts))
	cmd.AddCommand(newRemoveFromDeviceCmd(opts))
	cmd.AddCommand(newAssignDeviceGroupCmd(opts))
	cmd.AddCommand(newUnassignDeviceGroupCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all custom configuration profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/custom_configuration_profiles"
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
		Short: "Create a custom configuration profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			name, _ := cmd.Flags().GetString("name")
			values["name"] = name
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
			if cmd.Flags().Changed("reinstall-after-os-update") {
				ra, _ := cmd.Flags().GetBool("reinstall-after-os-update")
				if ra {
					values["reinstall_after_os_update"] = "true"
				} else {
					values["reinstall_after_os_update"] = "false"
				}
			}
			if cmd.Flags().Changed("declarative") {
				d, _ := cmd.Flags().GetBool("declarative")
				if d {
					values["declarative"] = "true"
				} else {
					values["declarative"] = "false"
				}
			}
			filePath, _ := cmd.Flags().GetString("mobileconfig")
			data, err := client.DoMultipart("POST", "/custom_configuration_profiles", values, "mobileconfig", filePath)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Profile name (required)")
	cmd.Flags().String("mobileconfig", "", "Path to .mobileconfig file (required)")
	cmd.Flags().String("user-scope", "", "User scope")
	cmd.Flags().Bool("attribute-support", false, "Enable attribute support")
	cmd.Flags().Bool("escape-attributes", false, "Escape attributes")
	cmd.Flags().Bool("reinstall-after-os-update", false, "Reinstall after OS update")
	cmd.Flags().Bool("declarative", false, "Declarative profile")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("mobileconfig")
	return cmd
}

func newDownloadCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <profile-id>",
		Short: "Download a custom configuration profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			body, filename, err := client.DoDownload("/custom_configuration_profiles/" + args[0] + "/download")
			if err != nil {
				return err
			}
			outPath, _ := cmd.Flags().GetString("output")
			if outPath == "" {
				if filename != "" {
					outPath = filename
				} else {
					outPath = "profile_" + args[0] + ".mobileconfig"
				}
			}
			if err := os.WriteFile(outPath, body, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("Profile downloaded to %s\n", outPath)
			return nil
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output file path")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <profile-id>",
		Short: "Update a custom configuration profile",
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
			if cmd.Flags().Changed("reinstall-after-os-update") {
				ra, _ := cmd.Flags().GetBool("reinstall-after-os-update")
				if ra {
					values["reinstall_after_os_update"] = "true"
				} else {
					values["reinstall_after_os_update"] = "false"
				}
			}
			if cmd.Flags().Changed("declarative") {
				d, _ := cmd.Flags().GetBool("declarative")
				if d {
					values["declarative"] = "true"
				} else {
					values["declarative"] = "false"
				}
			}
			filePath, _ := cmd.Flags().GetString("mobileconfig")
			if filePath != "" {
				data, err := client.DoMultipart("PATCH", "/custom_configuration_profiles/"+args[0], values, "mobileconfig", filePath)
				if err != nil {
					return err
				}
				return response.Print(opts.GetOutputFormat(), data, columns)
			}
			data, err := client.DoForm("PATCH", "/custom_configuration_profiles/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Profile name")
	cmd.Flags().String("mobileconfig", "", "Path to .mobileconfig file")
	cmd.Flags().String("user-scope", "", "User scope")
	cmd.Flags().Bool("attribute-support", false, "Enable attribute support")
	cmd.Flags().Bool("escape-attributes", false, "Escape attributes")
	cmd.Flags().Bool("reinstall-after-os-update", false, "Reinstall after OS update")
	cmd.Flags().Bool("declarative", false, "Declarative profile")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <profile-id>",
		Short: "Delete a custom configuration profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_configuration_profiles/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Custom configuration profile deleted.")
			return nil
		},
	}
}

func newPushToDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "push-to-device <profile-id> <device-id>",
		Short: "Push profile to device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/custom_configuration_profiles/"+args[0]+"/devices/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Profile pushed to device.")
			return nil
		},
	}
}

func newRemoveFromDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-from-device <profile-id> <device-id>",
		Short: "Remove profile from device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_configuration_profiles/" + args[0] + "/devices/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("Profile removed from device.")
			return nil
		},
	}
}

func newAssignDeviceGroupCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "assign-device-group <profile-id> <device-group-id>",
		Short: "Assign profile to device group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/custom_configuration_profiles/"+args[0]+"/device_groups/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Profile assigned to device group.")
			return nil
		},
	}
}

func newUnassignDeviceGroupCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-device-group <profile-id> <device-group-id>",
		Short: "Unassign profile from device group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/custom_configuration_profiles/" + args[0] + "/device_groups/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("Profile unassigned from device group.")
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
