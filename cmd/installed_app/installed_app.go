package installed_app

import (
	"fmt"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/response"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

var columns = []string{"id", "name", "identifier", "version", "managed"}

func NewCmdInstalledApp(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "installed-app",
		Short: "Manage installed apps",
	}
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newRequestManagementCmd(opts))
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <installed-app-id>",
		Short: "Get installed app details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/installed_apps/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <installed-app-id>",
		Short: "Delete an installed app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/installed_apps/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Installed app deleted.")
			return nil
		},
	}
}

func newUpdateCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "update <installed-app-id>",
		Short: "Update an installed app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/installed_apps/"+args[0]+"/update", nil)
			if err != nil {
				return err
			}
			fmt.Println("App update initiated.")
			return nil
		},
	}
}

func newRequestManagementCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "request-management <installed-app-id>",
		Short: "Request app management",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/installed_apps/"+args[0]+"/request_management", nil)
			if err != nil {
				return err
			}
			fmt.Println("Management requested.")
			return nil
		},
	}
}
