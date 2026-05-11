package profile

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

var columns = []string{"id", "name", "profile_identifier"}

func NewCmdProfile(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newAssignDeviceCmd(opts))
	cmd.AddCommand(newUnassignDeviceCmd(opts))
	cmd.AddCommand(newAssignDeviceGroupCmd(opts))
	cmd.AddCommand(newUnassignDeviceGroupCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/profiles"
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
	cmd.Flags().String("search", "", "Search by name or type")
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <profile-id>",
		Short: "Get profile details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/profiles/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newAssignDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "assign-device <profile-id> <device-id>",
		Short: "Assign profile to device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Post("/profiles/"+args[0]+"/devices/"+args[1], nil)
			if err != nil {
				return err
			}
			fmt.Println("Profile assigned to device.")
			return nil
		},
	}
}

func newUnassignDeviceCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "unassign-device <profile-id> <device-id>",
		Short: "Unassign profile from device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/profiles/" + args[0] + "/devices/" + args[1])
			if err != nil {
				return err
			}
			fmt.Println("Profile unassigned from device.")
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
			_, err = client.Post("/profiles/"+args[0]+"/device_groups/"+args[1], nil)
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
			_, err = client.Delete("/profiles/" + args[0] + "/device_groups/" + args[1])
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
