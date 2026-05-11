package assignment_group

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

var columns = []string{"id", "name", "auto_deploy"}

func NewCmdAssignmentGroup(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assignment-group",
		Aliases: []string{"ag"},
		Short:   "Manage assignment groups",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newAssignAppCmd(opts))
	cmd.AddCommand(newUnassignAppCmd(opts))
	cmd.AddCommand(newAssignDeviceCmd(opts))
	cmd.AddCommand(newUnassignDeviceCmd(opts))
	cmd.AddCommand(newAssignDeviceGroupCmd(opts))
	cmd.AddCommand(newUnassignDeviceGroupCmd(opts))
	cmd.AddCommand(newAssignProfileCmd(opts))
	cmd.AddCommand(newUnassignProfileCmd(opts))
	cmd.AddCommand(newPushAppsCmd(opts))
	cmd.AddCommand(newUpdateAppsCmd(opts))
	cmd.AddCommand(newSyncProfilesCmd(opts))
	cmd.AddCommand(newCloneCmd(opts))
	cmd.AddCommand(newCustomAttributesCmd(opts))
	cmd.AddCommand(newSetCustomAttributeCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all assignment groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/assignment_groups"
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
		Short: "Get assignment group details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/assignment_groups/" + args[0])
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
		Short: "Create an assignment group",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			name, _ := cmd.Flags().GetString("name")
			values["name"] = name
			if ad, _ := cmd.Flags().GetBool("auto-deploy"); ad {
				values["auto_deploy"] = "true"
			}
			if p, _ := cmd.Flags().GetString("priority"); p != "" {
				values["priority"] = p
			}
			if t, _ := cmd.Flags().GetString("type"); t != "" {
				values["type"] = t
			}
			if it, _ := cmd.Flags().GetString("install-type"); it != "" {
				values["install_type"] = it
			}
			if atl, _ := cmd.Flags().GetString("app-track-location"); atl != "" {
				values["app_track_location"] = atl
			}
			data, err := client.DoForm("POST", "/assignment_groups", values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Group name (required)")
	cmd.Flags().Bool("auto-deploy", false, "Auto deploy apps")
	cmd.Flags().String("priority", "", "Group priority")
	cmd.Flags().String("type", "", "Group type")
	cmd.Flags().String("install-type", "", "Install type")
	cmd.Flags().String("app-track-location", "", "App track location")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update an assignment group",
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
			if cmd.Flags().Changed("auto-deploy") {
				ad, _ := cmd.Flags().GetBool("auto-deploy")
				if ad {
					values["auto_deploy"] = "true"
				} else {
					values["auto_deploy"] = "false"
				}
			}
			if p, _ := cmd.Flags().GetString("priority"); p != "" {
				values["priority"] = p
			}
			data, err := client.DoForm("PATCH", "/assignment_groups/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("name", "", "Group name")
	cmd.Flags().Bool("auto-deploy", false, "Auto deploy apps")
	cmd.Flags().String("priority", "", "Group priority")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete an assignment group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/assignment_groups/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Assignment group deleted.")
			return nil
		},
	}
}

func newAssignAppCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-app <group-id> <app-id>",
		Short: "Assign app to group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if dt, _ := cmd.Flags().GetString("deployment-type"); dt != "" {
				values["deployment_type"] = dt
			}
			if it, _ := cmd.Flags().GetString("install-type"); it != "" {
				values["install_type"] = it
			}
			if len(values) > 0 {
				_, err = client.DoForm("POST", "/assignment_groups/"+args[0]+"/apps/"+args[1], values)
			} else {
				_, err = client.Post("/assignment_groups/"+args[0]+"/apps/"+args[1], nil)
			}
			if err != nil {
				return err
			}
			fmt.Println("App assigned.")
			return nil
		},
	}
	cmd.Flags().String("deployment-type", "", "Deployment type: standard or munki")
	cmd.Flags().String("install-type", "", "Install type: managed, self_serve, default_installs, managed_updates")
	return cmd
}

func newUnassignAppCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-app <group-id> <app-id>",
		Short: "Unassign app from group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgDelete(opts, func(a, b string) string { return "/assignment_groups/" + a + "/apps/" + b }, "App unassigned."),
	}
}

func newAssignDeviceCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-device <group-id> <device-id>",
		Short: "Assign device to group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if cmd.Flags().Changed("remove-others") {
				ro, _ := cmd.Flags().GetBool("remove-others")
				if ro {
					values["remove_others"] = "true"
				} else {
					values["remove_others"] = "false"
				}
			}
			if len(values) > 0 {
				_, err = client.DoForm("POST", "/assignment_groups/"+args[0]+"/devices/"+args[1], values)
			} else {
				_, err = client.Post("/assignment_groups/"+args[0]+"/devices/"+args[1], nil)
			}
			if err != nil {
				return err
			}
			fmt.Println("Device assigned.")
			return nil
		},
	}
	cmd.Flags().Bool("remove-others", false, "Remove device from other assignment groups")
	return cmd
}

func newUnassignDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-device <group-id> <device-id>",
		Short: "Unassign device from group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgDelete(opts, func(a, b string) string { return "/assignment_groups/" + a + "/devices/" + b }, "Device unassigned."),
	}
}

func newAssignDeviceGroupCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "assign-device-group <group-id> <device-group-id>",
		Short: "Assign device group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgPost(opts, func(a, b string) string { return "/assignment_groups/" + a + "/device_groups/" + b }, "Device group assigned."),
	}
}

func newUnassignDeviceGroupCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-device-group <group-id> <device-group-id>",
		Short: "Unassign device group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgDelete(opts, func(a, b string) string { return "/assignment_groups/" + a + "/device_groups/" + b }, "Device group unassigned."),
	}
}

func newAssignProfileCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "assign-profile <group-id> <profile-id>",
		Short: "Assign profile to group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgPost(opts, func(a, b string) string { return "/assignment_groups/" + a + "/profiles/" + b }, "Profile assigned."),
	}
}

func newUnassignProfileCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-profile <group-id> <profile-id>",
		Short: "Unassign profile from group",
		Args:  cobra.ExactArgs(2),
		RunE:  twoArgDelete(opts, func(a, b string) string { return "/assignment_groups/" + a + "/profiles/" + b }, "Profile unassigned."),
	}
}

func newPushAppsCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "push-apps <group-id>",
		Short: "Push apps to all devices in group",
		Args:  cobra.ExactArgs(1),
		RunE:  oneArgPost(opts, func(id string) string { return "/assignment_groups/" + id + "/push_apps" }, "Apps pushed."),
	}
}

func newUpdateAppsCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "update-apps <group-id>",
		Short: "Update apps on all devices in group",
		Args:  cobra.ExactArgs(1),
		RunE:  oneArgPost(opts, func(id string) string { return "/assignment_groups/" + id + "/update_apps" }, "Apps updated."),
	}
}

func newSyncProfilesCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "sync-profiles <group-id>",
		Short: "Sync profiles for group",
		Args:  cobra.ExactArgs(1),
		RunE:  oneArgPost(opts, func(id string) string { return "/assignment_groups/" + id + "/sync_profiles" }, "Profiles synced."),
	}
}

func newCloneCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clone <group-id>",
		Short: "Clone an assignment group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Post("/assignment_groups/"+args[0]+"/clone", nil)
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
			data, err := client.Get("/assignment_groups/" + args[0] + "/custom_attribute_values")
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
			_, err = client.DoForm("PUT", "/assignment_groups/"+args[0]+"/custom_attribute_values/"+args[1], map[string]string{"value": value})
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

// Helpers
func oneArgPost(opts *Options, pathFn func(string) string, msg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := opts.GetClient()
		if err != nil {
			return err
		}
		_, err = client.Post(pathFn(args[0]), nil)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	}
}

func twoArgPost(opts *Options, pathFn func(string, string) string, msg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := opts.GetClient()
		if err != nil {
			return err
		}
		_, err = client.Post(pathFn(args[0], args[1]), nil)
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	}
}

func twoArgDelete(opts *Options, pathFn func(string, string) string, msg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := opts.GetClient()
		if err != nil {
			return err
		}
		_, err = client.Delete(pathFn(args[0], args[1]))
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
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
