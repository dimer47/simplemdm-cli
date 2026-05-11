package device

import (
	"fmt"
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

var deviceColumns = []string{"id", "name", "serial_number", "model", "os_version", "status", "last_seen_at"}

func NewCmdDevice(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Manage devices",
	}

	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newInstalledAppsCmd(opts))
	cmd.AddCommand(newPushAppsCmd(opts))
	cmd.AddCommand(newRefreshCmd(opts))
	cmd.AddCommand(newLockCmd(opts))
	cmd.AddCommand(newWipeCmd(opts))
	cmd.AddCommand(newRestartCmd(opts))
	cmd.AddCommand(newShutdownCmd(opts))
	cmd.AddCommand(newClearPasscodeCmd(opts))
	cmd.AddCommand(newClearFirmwarePasswordCmd(opts))
	cmd.AddCommand(newUpdateOSCmd(opts))
	cmd.AddCommand(newBluetoothEnableCmd(opts))
	cmd.AddCommand(newBluetoothDisableCmd(opts))
	cmd.AddCommand(newRemoteDesktopEnableCmd(opts))
	cmd.AddCommand(newRemoteDesktopDisableCmd(opts))
	cmd.AddCommand(newRotateFirmwarePasswordCmd(opts))
	cmd.AddCommand(newClearRecoveryLockCmd(opts))
	cmd.AddCommand(newRotateRecoveryLockCmd(opts))
	cmd.AddCommand(newRotateFilevaultKeyCmd(opts))
	cmd.AddCommand(newSetAdminPasswordCmd(opts))
	cmd.AddCommand(newRotateAdminPasswordCmd(opts))
	cmd.AddCommand(newClearRestrictionsPasswordCmd(opts))
	cmd.AddCommand(newProfilesCmd(opts))
	cmd.AddCommand(newUsersCmd(opts))
	cmd.AddCommand(newDeleteUserCmd(opts))
	cmd.AddCommand(newSetTimeZoneCmd(opts))
	cmd.AddCommand(newUnenrollCmd(opts))
	cmd.AddCommand(newCustomAttributesCmd(opts))
	cmd.AddCommand(newSetCustomAttributeCmd(opts))
	cmd.AddCommand(newSetCustomAttributesCmd(opts))
	// Lost Mode
	cmd.AddCommand(newLostModeEnableCmd(opts))
	cmd.AddCommand(newLostModeDisableCmd(opts))
	cmd.AddCommand(newLostModePlaySoundCmd(opts))
	cmd.AddCommand(newLostModeUpdateLocationCmd(opts))

	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/devices"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction", "search", "include_awaiting_enrollment", "include_secret_custom_attributes")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return printListResponse(opts.GetOutputFormat(), data, deviceColumns)
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	cmd.Flags().String("search", "", "Search query")
	cmd.Flags().String("include-awaiting-enrollment", "", "Include awaiting enrollment: true or false")
	cmd.Flags().String("include-secret-custom-attributes", "", "Include secret custom attributes: true or false")
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get device details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/devices/" + args[0]
			params := buildQueryParams(cmd, "include_secret_custom_attributes")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, deviceColumns)
		},
	}
	cmd.Flags().String("include-secret-custom-attributes", "", "Include secret custom attributes: true or false")
	return cmd
}

func newCreateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a device",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			values := map[string]string{"name": name}
			if groupID, _ := cmd.Flags().GetString("group-id"); groupID != "" {
				values["group_id"] = groupID
			}
			if sgIDs, _ := cmd.Flags().GetString("static-group-ids"); sgIDs != "" {
				values["static_group_ids"] = sgIDs
			}
			data, err := client.DoForm("POST", "/devices", values)
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, deviceColumns)
		},
	}
	cmd.Flags().String("name", "", "Device name (required)")
	cmd.Flags().String("group-id", "", "Group ID (deprecated, use --static-group-ids)")
	cmd.Flags().String("static-group-ids", "", "Comma-separated static group IDs")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <device-id>",
		Short: "Update device name",
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
			if dn, _ := cmd.Flags().GetString("device-name"); dn != "" {
				values["device_name"] = dn
			}
			data, err := client.DoForm("PATCH", "/devices/"+args[0], values)
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, deviceColumns)
		},
	}
	cmd.Flags().String("name", "", "New SimpleMDM device name")
	cmd.Flags().String("device-name", "", "New device name (requires supervision, async)")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <device-id>",
		Short: "Delete a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/devices/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Device deleted.")
			return nil
		},
	}
}

func newInstalledAppsCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "installed-apps <device-id>",
		Short: "List installed apps on device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/devices/" + args[0] + "/installed_apps"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return printListResponse(opts.GetOutputFormat(), data, []string{"id", "name", "identifier", "version", "managed"})
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	return cmd
}

func newPushAppsCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "push-apps <device-id>",
		Short: "Push assigned apps to device",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/push_apps" }, "Apps pushed."),
	}
}

func newRefreshCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh <device-id>",
		Short: "Refresh device information",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/refresh" }, "Device refresh requested."),
	}
}

func newLockCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock <device-id>",
		Short: "Lock a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if msg, _ := cmd.Flags().GetString("message"); msg != "" {
				values["message"] = msg
			}
			if phone, _ := cmd.Flags().GetString("phone-number"); phone != "" {
				values["phone_number"] = phone
			}
			if pin, _ := cmd.Flags().GetString("pin"); pin != "" {
				values["pin"] = pin
			}
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/lock", values)
			if err != nil {
				return err
			}
			fmt.Println("Device locked.")
			return nil
		},
	}
	cmd.Flags().String("message", "", "Lock message")
	cmd.Flags().String("phone-number", "", "Phone number to display")
	cmd.Flags().String("pin", "", "PIN code")
	return cmd
}

func newWipeCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wipe <device-id>",
		Short: "Wipe a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if pin, _ := cmd.Flags().GetString("pin"); pin != "" {
				values["pin"] = pin
			}
			if cmd.Flags().Changed("clear-custom-attributes") {
				values["clear_custom_attributes"] = "true"
			}
			if cmd.Flags().Changed("disable-activation-lock") {
				values["disable_activation_lock"] = "true"
			}
			if cmd.Flags().Changed("preserve-data-plan") {
				values["preserve_data_plan"] = "true"
			}
			if cmd.Flags().Changed("disallow-proximity-setup") {
				values["disallow_proximity_setup"] = "true"
			}
			if cmd.Flags().Changed("return-to-service") {
				values["return_to_service"] = "true"
			}
			if wid, _ := cmd.Flags().GetString("wifi-network-id"); wid != "" {
				values["wifi_network_id"] = wid
			}
			if ob, _ := cmd.Flags().GetString("obliteration-behavior"); ob != "" {
				values["obliteration_behavior"] = ob
			}
			if cmd.Flags().Changed("unassign-direct-profiles") {
				values["unassign_direct_profiles"] = "true"
			}
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/wipe", values)
			if err != nil {
				return err
			}
			fmt.Println("Device wipe initiated.")
			return nil
		},
	}
	cmd.Flags().String("pin", "", "PIN code")
	cmd.Flags().Bool("clear-custom-attributes", false, "Clear custom attributes")
	cmd.Flags().Bool("disable-activation-lock", false, "Disable activation lock")
	cmd.Flags().Bool("preserve-data-plan", false, "Preserve data plan (iOS)")
	cmd.Flags().Bool("disallow-proximity-setup", false, "Disallow proximity setup (iOS)")
	cmd.Flags().Bool("return-to-service", false, "Return to service (iOS 17+/tvOS 18+)")
	cmd.Flags().String("wifi-network-id", "", "WiFi network ID for return to service")
	cmd.Flags().String("obliteration-behavior", "", "Obliteration behavior (macOS 12+)")
	cmd.Flags().Bool("unassign-direct-profiles", false, "Unassign direct profiles")
	return cmd
}

func newRestartCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart <device-id>",
		Short: "Restart a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if cmd.Flags().Changed("rebuild-kernel-cache") {
				rkc, _ := cmd.Flags().GetBool("rebuild-kernel-cache")
				if rkc {
					values["rebuild_kernel_cache"] = "true"
				}
			}
			if cmd.Flags().Changed("notify-user") {
				nu, _ := cmd.Flags().GetBool("notify-user")
				if nu {
					values["notify_user"] = "true"
				}
			}
			if len(values) > 0 {
				_, err = client.DoForm("POST", "/devices/"+args[0]+"/restart", values)
			} else {
				_, err = client.Post("/devices/"+args[0]+"/restart", nil)
			}
			if err != nil {
				return err
			}
			fmt.Println("Device restart requested.")
			return nil
		},
	}
	cmd.Flags().Bool("rebuild-kernel-cache", false, "Rebuild kernel cache (macOS 11+)")
	cmd.Flags().Bool("notify-user", false, "Notify user (macOS 11.3+)")
	return cmd
}

func newShutdownCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "shutdown <device-id>",
		Short: "Shut down a device",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/shutdown" }, "Device shutdown requested."),
	}
}

func newClearPasscodeCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-passcode <device-id>",
		Short: "Clear device passcode",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/clear_passcode" }, "Passcode cleared."),
	}
}

func newClearFirmwarePasswordCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-firmware-password <device-id>",
		Short: "Clear firmware password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/clear_firmware_password" }, "Firmware password cleared."),
	}
}

func newUpdateOSCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-os <device-id>",
		Short: "Update device OS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if mode, _ := cmd.Flags().GetString("os-update-mode"); mode != "" {
				values["os_update_mode"] = mode
			}
			if vt, _ := cmd.Flags().GetString("version-type"); vt != "" {
				values["version_type"] = vt
			}
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/update_os", values)
			if err != nil {
				return err
			}
			fmt.Println("OS update initiated.")
			return nil
		},
	}
	cmd.Flags().String("os-update-mode", "", "Update mode: smart_update, download_only, notify_only, install_asap, force_update")
	cmd.Flags().String("version-type", "", "Version type: latest_minor_version or latest_major_version")
	return cmd
}

func newBluetoothEnableCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "bluetooth-enable <device-id>",
		Short: "Enable Bluetooth",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/bluetooth" }, "Bluetooth enabled."),
	}
}

func newBluetoothDisableCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "bluetooth-disable <device-id>",
		Short: "Disable Bluetooth",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/devices/" + args[0] + "/bluetooth")
			if err != nil {
				return err
			}
			fmt.Println("Bluetooth disabled.")
			return nil
		},
	}
}

func newRemoteDesktopEnableCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "remote-desktop-enable <device-id>",
		Short: "Enable Remote Desktop",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/remote_desktop" }, "Remote Desktop enabled."),
	}
}

func newRemoteDesktopDisableCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "remote-desktop-disable <device-id>",
		Short: "Disable Remote Desktop",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/devices/" + args[0] + "/remote_desktop")
			if err != nil {
				return err
			}
			fmt.Println("Remote Desktop disabled.")
			return nil
		},
	}
}

func newRotateFirmwarePasswordCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate-firmware-password <device-id>",
		Short: "Rotate firmware password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/rotate_firmware_password" }, "Firmware password rotated."),
	}
}

func newClearRecoveryLockCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-recovery-lock <device-id>",
		Short: "Clear recovery lock password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/clear_recovery_lock_password" }, "Recovery lock password cleared."),
	}
}

func newRotateRecoveryLockCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate-recovery-lock <device-id>",
		Short: "Rotate recovery lock password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/rotate_recovery_lock_password" }, "Recovery lock password rotated."),
	}
}

func newRotateFilevaultKeyCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate-filevault-key <device-id>",
		Short: "Rotate FileVault key",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/rotate_filevault_key" }, "FileVault key rotated."),
	}
}

func newSetAdminPasswordCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-admin-password <device-id>",
		Short: "Set admin password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			pw, _ := cmd.Flags().GetString("new-password")
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/set_admin_password", map[string]string{"new_password": pw})
			if err != nil {
				return err
			}
			fmt.Println("Admin password set.")
			return nil
		},
	}
	cmd.Flags().String("new-password", "", "New admin password (required)")
	_ = cmd.MarkFlagRequired("new-password")
	return cmd
}

func newRotateAdminPasswordCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate-admin-password <device-id>",
		Short: "Rotate admin password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/rotate_admin_password" }, "Admin password rotated."),
	}
}

func newClearRestrictionsPasswordCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-restrictions-password <device-id>",
		Short: "Clear restrictions password",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/clear_restrictions_password" }, "Restrictions password cleared."),
	}
}

func newProfilesCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles <device-id>",
		Short: "List device profiles",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/devices/" + args[0] + "/profiles"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return printListResponse(opts.GetOutputFormat(), data, []string{"id", "name", "profile_identifier"})
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	return cmd
}

func newUsersCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users <device-id>",
		Short: "List device users (macOS only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/devices/" + args[0] + "/users"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return printListResponse(opts.GetOutputFormat(), data, nil)
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	return cmd
}

func newDeleteUserCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-user <device-id> <user-id>",
		Short: "Delete a user from device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/devices/" + args[0] + "/users/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("User deleted from device.")
			return nil
		},
	}
}

func newSetTimeZoneCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-timezone <device-id>",
		Short: "Set device time zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tz, _ := cmd.Flags().GetString("timezone")
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/set_time_zone", map[string]string{"time_zone": tz})
			if err != nil {
				return err
			}
			fmt.Println("Time zone set.")
			return nil
		},
	}
	cmd.Flags().String("timezone", "", "Time zone (required)")
	_ = cmd.MarkFlagRequired("timezone")
	return cmd
}

func newUnenrollCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unenroll <device-id>",
		Short: "Unenroll a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if cmd.Flags().Changed("unassign-direct-profiles") {
				values["unassign_direct_profiles"] = "true"
			}
			if len(values) > 0 {
				_, err = client.DoForm("POST", "/devices/"+args[0]+"/unenroll", values)
			} else {
				_, err = client.Post("/devices/"+args[0]+"/unenroll", nil)
			}
			if err != nil {
				return err
			}
			fmt.Println("Device unenrolled.")
			return nil
		},
	}
	cmd.Flags().Bool("unassign-direct-profiles", false, "Unassign direct profiles")
	return cmd
}

func newCustomAttributesCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "custom-attributes <device-id>",
		Short: "List custom attribute values for device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/devices/" + args[0] + "/custom_attribute_values")
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newSetCustomAttributeCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-custom-attribute <device-id> <attribute-name>",
		Short: "Set a custom attribute value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			value, _ := cmd.Flags().GetString("value")
			_, err = client.DoForm("PUT", "/devices/"+args[0]+"/custom_attribute_values/"+args[1], map[string]string{"value": value})
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

func newSetCustomAttributesCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-custom-attributes <device-id>",
		Short: "Set multiple custom attribute values (JSON)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			jsonData, _ := cmd.Flags().GetString("data")
			data, err := client.Do("PUT", "/devices/"+args[0]+"/custom_attribute_values", strings.NewReader(jsonData))
			if err != nil {
				return err
			}
			return printResponse(opts.GetOutputFormat(), data, nil)
		},
	}
	cmd.Flags().String("data", "", "JSON key-value pairs (required)")
	_ = cmd.MarkFlagRequired("data")
	return cmd
}

// Lost Mode commands
func newLostModeEnableCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lost-mode-enable <device-id>",
		Short: "Enable Lost Mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if msg, _ := cmd.Flags().GetString("message"); msg != "" {
				values["message"] = msg
			}
			if phone, _ := cmd.Flags().GetString("phone-number"); phone != "" {
				values["phone_number"] = phone
			}
			if fn, _ := cmd.Flags().GetString("footnote"); fn != "" {
				values["footnote"] = fn
			}
			_, err = client.DoForm("POST", "/devices/"+args[0]+"/lost_mode", values)
			if err != nil {
				return err
			}
			fmt.Println("Lost Mode enabled.")
			return nil
		},
	}
	cmd.Flags().String("message", "", "Message to display")
	cmd.Flags().String("phone-number", "", "Phone number to display")
	cmd.Flags().String("footnote", "", "Footnote to display")
	return cmd
}

func newLostModeDisableCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "lost-mode-disable <device-id>",
		Short: "Disable Lost Mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/devices/" + args[0] + "/lost_mode")
			if err != nil {
				return err
			}
			fmt.Println("Lost Mode disabled.")
			return nil
		},
	}
}

func newLostModePlaySoundCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "lost-mode-play-sound <device-id>",
		Short: "Play Lost Mode sound",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/lost_mode/play_sound" }, "Lost Mode sound played."),
	}
}

func newLostModeUpdateLocationCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "lost-mode-update-location <device-id>",
		Short: "Request location update",
		Args:  cobra.ExactArgs(1),
		RunE:  simplePostAction(opts, func(id string) string { return "/devices/" + id + "/lost_mode/update_location" }, "Location update requested."),
	}
}

// Helper functions
func simplePostAction(opts *Options, pathFn func(string) string, msg string) func(cmd *cobra.Command, args []string) error {
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

func buildQueryParams(cmd *cobra.Command, params ...string) string {
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
