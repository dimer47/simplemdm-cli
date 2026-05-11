package log

import (
	"strings"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/response"
	"github.com/spf13/cobra"
)

type Options struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

var columns = []string{"id", "namespace", "event_type", "created_at"}

func NewCmdLog(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "View logs",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/logs"
			params := buildQP(cmd, "limit", "starting_after", "direction", "serial_number")
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
	cmd.Flags().String("serial-number", "", "Filter by device serial number")
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <log-id>",
		Short: "Get log details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/logs/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
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
