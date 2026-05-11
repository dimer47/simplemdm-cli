package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dimer47/simplemdm-cli/cmd/account"
	"github.com/dimer47/simplemdm-cli/cmd/app"
	"github.com/dimer47/simplemdm-cli/cmd/assignment_group"
	"github.com/dimer47/simplemdm-cli/cmd/auth"
	"github.com/dimer47/simplemdm-cli/cmd/custom_attribute"
	"github.com/dimer47/simplemdm-cli/cmd/custom_configuration_profile"
	"github.com/dimer47/simplemdm-cli/cmd/custom_declaration"
	"github.com/dimer47/simplemdm-cli/cmd/dep_server"
	"github.com/dimer47/simplemdm-cli/cmd/device"
	"github.com/dimer47/simplemdm-cli/cmd/device_group"
	"github.com/dimer47/simplemdm-cli/cmd/enrollment"
	"github.com/dimer47/simplemdm-cli/cmd/installed_app"
	cmdlog "github.com/dimer47/simplemdm-cli/cmd/log"
	"github.com/dimer47/simplemdm-cli/cmd/mcp"
	"github.com/dimer47/simplemdm-cli/cmd/profile"
	"github.com/dimer47/simplemdm-cli/cmd/push_certificate"
	"github.com/dimer47/simplemdm-cli/cmd/script"
	"github.com/dimer47/simplemdm-cli/cmd/script_job"
	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/config"
	"github.com/dimer47/simplemdm-cli/internal/keychain"
	"github.com/dimer47/simplemdm-cli/internal/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "simplemdm-cli",
	Short: "CLI for SimpleMDM - Manage iOS/macOS devices",
	Long:  `A command-line interface for SimpleMDM to manage your Apple devices, apps, profiles, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "[DEBUG] simplemdm-cli %s (%s) built %s\n", Version, Commit, Date)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("api-key", "k", "", "SimpleMDM API key (env: SMDM_API_KEY)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml, csv (env: SMDM_OUTPUT)")
	rootCmd.PersistentFlags().Bool("json", false, "Shortcut for --output json")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode (env: SMDM_DEBUG)")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output (env: NO_COLOR)")
	rootCmd.PersistentFlags().StringP("context", "c", "", "Active context (env: SMDM_CONTEXT)")

	_ = viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	_ = viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))

	_ = viper.BindEnv("api-key", "SMDM_API_KEY")
	_ = viper.BindEnv("output", "SMDM_OUTPUT")
	_ = viper.BindEnv("debug", "SMDM_DEBUG")
	_ = viper.BindEnv("context", "SMDM_CONTEXT")

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("simplemdm-cli %s (%s) built %s\n", Version, Commit, Date)
		},
	})

	// Auth
	rootCmd.AddCommand(auth.NewCmdAuth())

	// MCP
	rootCmd.AddCommand(mcp.NewCmdMCP(&mcp.MCPOptions{
		GetClient: getClient,
	}))

	// Domain commands
	opts := &CmdOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
	}

	rootCmd.AddCommand(account.NewCmdAccount(&account.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(app.NewCmdApp(&app.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(assignment_group.NewCmdAssignmentGroup(&assignment_group.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(custom_attribute.NewCmdCustomAttribute(&custom_attribute.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(custom_configuration_profile.NewCmdCustomConfigurationProfile(&custom_configuration_profile.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(custom_declaration.NewCmdCustomDeclaration(&custom_declaration.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(dep_server.NewCmdDepServer(&dep_server.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(device_group.NewCmdDeviceGroup(&device_group.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(device.NewCmdDevice(&device.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(enrollment.NewCmdEnrollment(&enrollment.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(installed_app.NewCmdInstalledApp(&installed_app.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(cmdlog.NewCmdLog(&cmdlog.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(profile.NewCmdProfile(&profile.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(push_certificate.NewCmdPushCertificate(&push_certificate.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(script.NewCmdScript(&script.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))
	rootCmd.AddCommand(script_job.NewCmdScriptJob(&script_job.Options{GetClient: opts.GetClient, GetOutputFormat: opts.GetOutputFormat}))

	// Self-update
	rootCmd.AddCommand(newSelfUpdateCmd())

	// Completion
	rootCmd.AddCommand(newCompletionCmd())
}

type CmdOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
}

func initConfig() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
}

func getClient() (*api.Client, error) {
	// Priority: flag > env > keychain
	apiKey := viper.GetString("api-key")

	if apiKey == "" {
		// Try keychain with context
		contextName := viper.GetString("context")
		if contextName == "" {
			cfg, err := config.Load()
			if err == nil {
				contextName = cfg.DefaultContext
			}
		}
		if contextName != "" {
			if key, err := keychain.Get(contextName); err == nil && key != "" {
				apiKey = key
			}
		}
	}

	if apiKey == "" {
		return nil, fmt.Errorf("no API key configured. Use 'simplemdm-cli auth login' or set SMDM_API_KEY")
	}

	client := api.NewClient(apiKey)
	client.SetDebug(viper.GetBool("debug"))
	return client, nil
}

func getOutputFormat() string {
	if viper.GetBool("json") {
		return "json"
	}
	return viper.GetString("output")
}

func Execute() {
	// Launch background update check (non-blocking, 1s timeout)
	type updateMsg struct {
		result *update.CheckResult
	}
	updateCh := make(chan updateMsg, 1)
	go func() {
		result := update.Check(Version)
		updateCh <- updateMsg{result: result}
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	// After command execution, check if an update notification is available
	select {
	case msg := <-updateCh:
		if msg.result != nil && msg.result.UpdateAvailable && msg.result.LatestVersion != "" {
			fmt.Fprintf(os.Stderr, "\nA new version of simplemdm-cli is available: %s -> %s\n", msg.result.CurrentVersion, msg.result.LatestVersion)
			fmt.Fprintf(os.Stderr, "Run 'simplemdm-cli self-update' to update.\n")
		}
	case <-time.After(1 * time.Second):
		// Timeout: don't block the CLI
	}
}
