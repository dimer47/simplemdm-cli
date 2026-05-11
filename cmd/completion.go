package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long:  `Generate autocompletion scripts for bash, zsh, fish, or powershell.`,
	}

	completionCmd.AddCommand(&cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		Long: `Generate the autocompletion script for bash.

To load completions in your current shell session:

  source <(simplemdm-cli completion bash)

To load completions for every new session, execute once:

  # Linux:
  simplemdm-cli completion bash > /etc/bash_completion.d/simplemdm-cli

  # macOS:
  simplemdm-cli completion bash > $(brew --prefix)/etc/bash_completion.d/simplemdm-cli
`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		Long: `Generate the autocompletion script for zsh.

To load completions in your current shell session:

  source <(simplemdm-cli completion zsh)

To load completions for every new session, execute once:

  simplemdm-cli completion zsh > "${fpath[1]}/_simplemdm-cli"

You will need to start a new shell for this setup to take effect.
`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenZshCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion script",
		Long: `Generate the autocompletion script for fish.

To load completions in your current shell session:

  simplemdm-cli completion fish | source

To load completions for every new session, execute once:

  simplemdm-cli completion fish > ~/.config/fish/completions/simplemdm-cli.fish
`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(os.Stdout, true)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "powershell",
		Short: "Generate powershell completion script",
		Long: `Generate the autocompletion script for powershell.

To load completions in your current shell session:

  simplemdm-cli completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.
`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenPowerShellCompletion(os.Stdout)
		},
	})

	return completionCmd
}
