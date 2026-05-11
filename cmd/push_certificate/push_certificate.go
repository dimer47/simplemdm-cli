package push_certificate

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

var columns = []string{"apple_id", "expires_at"}

func NewCmdPushCertificate(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push-certificate",
		Short: "Manage push certificate",
	}
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newGetSCSRCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	return cmd
}

func newGetCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get push certificate details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/push_certificate")
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
}

func newUpdateCmd(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update push certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			values := make(map[string]string)
			if appleID, _ := cmd.Flags().GetString("apple-id"); appleID != "" {
				values["apple_id"] = appleID
			}
			filePath, _ := cmd.Flags().GetString("file")
			data, err := client.DoMultipart("PUT", "/push_certificate", values, "file", filePath)
			if err != nil {
				return err
			}
			fmt.Println("Push certificate updated.")
			return response.Print(opts.GetOutputFormat(), data, columns)
		},
	}
	cmd.Flags().String("file", "", "Path to Apple push certificate file (required)")
	cmd.Flags().String("apple-id", "", "Apple ID")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newGetSCSRCmd(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "scsr",
		Short: "Get signed CSR",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			data, err := client.Get("/push_certificate/scsr")
			if err != nil {
				return err
			}
			return response.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}
