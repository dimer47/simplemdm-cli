package enrollment

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

var columns = []string{"id", "url"}

func NewCmdEnrollment(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enrollment",
		Short: "Manage enrollments",
	}
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newSendInvitationCmd(opts))
	return cmd
}

func newListCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all enrollments",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			path := "/enrollments"
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
		Use:   "get <enrollment-id>",
		Short: "Get enrollment details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/enrollments/" + args[0])
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newDeleteCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <enrollment-id>",
		Short: "Delete an enrollment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			_, err = client.Delete("/enrollments/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Enrollment deleted.")
			return nil
		},
	}
}

func newSendInvitationCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-invitation <enrollment-id>",
		Short: "Send enrollment invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			contact, _ := cmd.Flags().GetString("contact")
			_, err = client.DoForm("POST", "/enrollments/"+args[0]+"/invitations", map[string]string{"contact": contact})
			if err != nil {
				return err
			}
			fmt.Println("Invitation sent.")
			return nil
		},
	}
	cmd.Flags().String("contact", "", "Email or phone number (required)")
	_ = cmd.MarkFlagRequired("contact")
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
