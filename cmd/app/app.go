package app

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

var appColumns = []string{"id", "name", "app_type", "bundle_identifier", "version"}

func NewCmdApp(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Manage apps",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newInstallsCmd(opts))
	cmd.AddCommand(newManagedConfigsListCmd(opts))
	cmd.AddCommand(newManagedConfigsCreateCmd(opts))
	cmd.AddCommand(newManagedConfigsPushCmd(opts))
	cmd.AddCommand(newManagedConfigsDeleteCmd(opts))
	cmd.AddCommand(newMunkiPkginfoUpdateCmd(opts))
	cmd.AddCommand(newMunkiPkginfoDeleteCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/apps"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction", "include_shared")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, appColumns)
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	cmd.Flags().String("include-shared", "", "Include shared apps: true or false")
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <app-id>",
		Short: "Get app details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/apps/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, appColumns)
		},
	}
}

func newCreateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an app",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if name, _ := cmd.Flags().GetString("name"); name != "" {
				values["name"] = name
			}
			if appStoreID, _ := cmd.Flags().GetString("app-store-id"); appStoreID != "" {
				values["app_store_id"] = appStoreID
			}
			if bundleID, _ := cmd.Flags().GetString("bundle-id"); bundleID != "" {
				values["bundle_id"] = bundleID
			}
			filePath, _ := cmd.Flags().GetString("binary")
			if filePath != "" {
				data, err := client.DoMultipart("POST", "/apps", values, "binary", filePath)
				if err != nil {
					return err
				}
				return response.Print(opts.GetOutputFormat(), data, appColumns)
			}
			data, err := client.DoForm("POST", "/apps", values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, appColumns)
		},
	}
	cmd.Flags().String("name", "", "App name")
	cmd.Flags().String("app-store-id", "", "App Store ID")
	cmd.Flags().String("bundle-id", "", "Bundle ID")
	cmd.Flags().String("binary", "", "Path to app binary file")
	return cmd
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <app-id>",
		Short: "Update an app",
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
			if deployTo, _ := cmd.Flags().GetString("deploy-to"); deployTo != "" {
				values["deploy_to"] = deployTo
			}
			filePath, _ := cmd.Flags().GetString("binary")
			if filePath != "" {
				data, err := client.DoMultipart("PATCH", "/apps/"+args[0], values, "binary", filePath)
				if err != nil {
					return err
				}
				return response.Print(opts.GetOutputFormat(), data, appColumns)
			}
			data, err := client.DoForm("PATCH", "/apps/"+args[0], values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, appColumns)
		},
	}
	cmd.Flags().String("name", "", "New app name")
	cmd.Flags().String("binary", "", "Path to app binary file")
	cmd.Flags().String("deploy-to", "", "Deploy to: none, outdated, or all")
	return cmd
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <app-id>",
		Short: "Delete an app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/apps/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("App deleted.")
			return nil
		},
	}
}

func newInstallsCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "installs <app-id>",
		Short: "List app installs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/apps/" + args[0] + "/installs"
			params := buildQueryParams(cmd, "limit", "starting_after", "direction")
			if params != "" {
				path += "?" + params
			}
			data, err := client.Get(path)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, nil)
		},
	}
	cmd.Flags().Int("limit", 0, "Limit results")
	cmd.Flags().String("starting-after", "", "Cursor for pagination")
	cmd.Flags().String("direction", "", "Sort direction: asc or desc")
	return cmd
}

func newManagedConfigsListCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "managed-configs <app-id>",
		Short: "List managed app configs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/apps/" + args[0] + "/managed_configs")
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, []string{"id", "key", "value", "value_type"})
		},
	}
}

func newManagedConfigsCreateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "managed-config-create <app-id>",
		Short: "Create a managed config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			key, _ := cmd.Flags().GetString("key")
			value, _ := cmd.Flags().GetString("value")
			valueType, _ := cmd.Flags().GetString("value-type")
			values := map[string]string{"key": key, "value": value, "value_type": valueType}
			data, err := client.DoMultipart("POST", "/apps/"+args[0]+"/managed_configs", values, "", "")
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, []string{"id", "key", "value", "value_type"})
		},
	}
	cmd.Flags().String("key", "", "Config key (required)")
	cmd.Flags().String("value", "", "Config value (required)")
	cmd.Flags().String("value-type", "string", "Value type: boolean, date, float, float array, integer, integer array, string, string array")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newManagedConfigsPushCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "managed-configs-push <app-id>",
		Short: "Push managed configs to devices",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/apps/"+args[0]+"/managed_configs/push", nil)
			if err != nil {
				return err
			}
			fmt.Println("Managed configs pushed.")
			return nil
		},
	}
}

func newManagedConfigsDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "managed-config-delete <app-id> <config-id>",
		Short: "Delete a managed config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/apps/" + args[0] + "/managed_configs/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("Managed config deleted.")
			return nil
		},
	}
}

func newMunkiPkginfoUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "munki-pkginfo-update <app-id>",
		Short: "Upload Munki pkginfo file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			filePath, _ := cmd.Flags().GetString("file")
			_, err = client.DoMultipart("POST", "/apps/"+args[0]+"/munki_pkginfo", nil, "file", filePath)
			if err != nil {
				return err
			}
			fmt.Println("Munki pkginfo uploaded.")
			return nil
		},
	}
	cmd.Flags().String("file", "", "Path to pkginfo XML/PLIST file (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newMunkiPkginfoDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "munki-pkginfo-delete <app-id>",
		Short: "Delete Munki pkginfo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/apps/" + args[0] + "/munki_pkginfo")
			if err != nil {
				return err
			}
			fmt.Println("Munki pkginfo deleted.")
			return nil
		},
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
