package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/dimer47/simplemdm-cli/internal/api"
	"github.com/dimer47/simplemdm-cli/internal/config"
	"github.com/dimer47/simplemdm-cli/internal/keychain"
	"github.com/dimer47/simplemdm-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication contexts",
	}

	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newSwitchCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newRemoveCmd())

	return cmd
}

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Add a new authentication context",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !keychain.IsAvailable() {
				return fmt.Errorf("system keychain is not available")
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Context name (default): ")
			contextName, _ := reader.ReadString('\n')
			contextName = strings.TrimSpace(contextName)
			if contextName == "" {
				contextName = "default"
			}

			fmt.Print("SimpleMDM API key: ")
			apiKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read API key: %w", err)
			}
			fmt.Println()
			apiKey := strings.TrimSpace(string(apiKeyBytes))

			if apiKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}

			// Validate the API key
			client := api.NewClient(apiKey)
			_, err = client.Get("/account")
			if err != nil {
				return fmt.Errorf("invalid API key: %w", err)
			}

			// Store in keychain
			if err := keychain.Set(contextName, apiKey); err != nil {
				return fmt.Errorf("failed to store API key in keychain: %w", err)
			}

			// Save context in config
			cfg, err := config.Load()
			if err != nil {
				cfg = config.NewConfig()
			}

			cfg.SetContext(contextName, config.Context{Name: contextName})
			cfg.DefaultContext = contextName

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Context '%s' configured and set as default.\n", contextName)
			fmt.Printf("API key stored securely in system keychain.\n")
			return nil
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.DefaultContext == "" {
				fmt.Println("No context configured. Run 'simplemdm-cli auth login' to configure.")
				return nil
			}

			fmt.Printf("Active context: %s\n", cfg.DefaultContext)

			apiKey, err := keychain.Get(cfg.DefaultContext)
			if err != nil {
				fmt.Printf("API key: not found in keychain\n")
				fmt.Printf("Status: not authenticated\n")
				return nil
			}

			fmt.Printf("API key: %s\n", output.MaskToken(apiKey))

			// Test the key
			client := api.NewClient(apiKey)
			data, err := client.Get("/account")
			if err != nil {
				fmt.Printf("Status: invalid (%s)\n", err)
				return nil
			}

			fmt.Printf("Status: authenticated\n")
			fmt.Printf("Account: %s\n", strings.TrimSpace(string(data)))
			return nil
		},
	}
}

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <context>",
		Short: "Switch to a different context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if _, ok := cfg.GetContext(contextName); !ok {
				return fmt.Errorf("context '%s' not found", contextName)
			}

			cfg.DefaultContext = contextName
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Switched to context '%s'.\n", contextName)
			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if len(cfg.Contexts) == 0 {
				fmt.Println("No contexts configured. Run 'simplemdm-cli auth login' to configure.")
				return nil
			}

			for name := range cfg.Contexts {
				marker := "  "
				if name == cfg.DefaultContext {
					marker = "* "
				}
				tokenStatus := "no key"
				if key, err := keychain.Get(name); err == nil && key != "" {
					tokenStatus = output.MaskToken(key)
				}
				fmt.Printf("%s%s (%s)\n", marker, name, tokenStatus)
			}
			return nil
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <context>",
		Short: "Remove a context and its API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if _, ok := cfg.GetContext(contextName); !ok {
				return fmt.Errorf("context '%s' not found", contextName)
			}

			// Delete from keychain
			_ = keychain.Delete(contextName)

			// Delete from config
			cfg.DeleteContext(contextName)
			if cfg.DefaultContext == contextName {
				cfg.DefaultContext = ""
				// Set first remaining context as default
				for name := range cfg.Contexts {
					cfg.DefaultContext = name
					break
				}
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Context '%s' removed.\n", contextName)
			return nil
		},
	}
}
