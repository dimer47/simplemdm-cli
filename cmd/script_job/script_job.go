package script_job

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

var columns = []string{"id", "script_id", "status"}

func NewCmdScriptJob(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script-job",
		Short: "Manage script jobs",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newCancelCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all script jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/script_jobs"
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
		Use:   "get <job-id>",
		Short: "Get script job details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/script_jobs/" + args[0])
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
		Short: "Create a script job",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			scriptID, _ := cmd.Flags().GetString("script-id")
			values["script_id"] = scriptID
			if deviceIDs, _ := cmd.Flags().GetString("device-ids"); deviceIDs != "" {
				values["device_ids"] = deviceIDs
			}
			if agIDs, _ := cmd.Flags().GetString("assignment-group-ids"); agIDs != "" {
				values["assignment_group_ids"] = agIDs
			}
			if gIDs, _ := cmd.Flags().GetString("group-ids"); gIDs != "" {
				values["group_ids"] = gIDs
			}
			if ca, _ := cmd.Flags().GetString("custom-attribute"); ca != "" {
				values["custom_attribute"] = ca
			}
			if car, _ := cmd.Flags().GetString("custom-attribute-regex"); car != "" {
				values["custom_attribute_regex"] = car
			}
			data, err := client.DoForm("POST", "/script_jobs", values)
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("script-id", "", "Script ID (required)")
	cmd.Flags().String("device-ids", "", "Comma-separated device IDs")
	cmd.Flags().String("assignment-group-ids", "", "Comma-separated assignment group IDs")
	cmd.Flags().String("group-ids", "", "Comma-separated device group IDs (deprecated)")
	cmd.Flags().String("custom-attribute", "", "Custom attribute filter")
	cmd.Flags().String("custom-attribute-regex", "", "Custom attribute regex filter")
	_ = cmd.MarkFlagRequired("script-id")
	return cmd
}

func newCancelCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <job-id>",
		Short: "Cancel a script job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/script_jobs/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Script job cancelled.")
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
