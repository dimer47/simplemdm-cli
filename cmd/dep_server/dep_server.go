package dep_server

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

var columns = []string{"id", "server_name", "organization_name"}

func NewCmdDepServer(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dep-server",
		Aliases: []string{"dep"},
		Short:   "Manage DEP servers",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newDevicesCmd(opts))
	cmd.AddCommand(newDeviceGetCmd(opts))
	cmd.AddCommand(newSyncCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all DEP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/dep_servers"
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
		Use:   "get <server-id>",
		Short: "Get DEP server details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/dep_servers/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newDevicesCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices <server-id>",
		Short: "List DEP devices for server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/dep_servers/" + args[0] + "/dep_devices"
			params := buildQP(cmd, "limit", "starting_after", "direction")
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

func newDeviceGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "device-get <server-id> <dep-device-id>",
		Short: "Get a specific DEP device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/dep_servers/" + args[0] + "/dep_devices/" + args[1])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newSyncCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "sync <server-id>",
		Short: "Sync DEP server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/dep_servers/"+args[0]+"/sync", nil)
			if err != nil {
				return err
			}
			fmt.Println("DEP server sync initiated.")
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
