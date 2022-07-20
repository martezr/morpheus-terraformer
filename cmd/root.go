package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdRoot defines the command line arguments
func NewCmdRoot() *cobra.Command {
	var cmd = &cobra.Command{
		//		SilenceUsage: true,
		Use:  "morpheus-terraformer",
		Long: `Generate Terraform code from existing Morpheus resources`,
	}
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(generateCmd)
	//cmd.AddCommand(installProviderCmd)
	return cmd
}

// Execute runs the commands
func Execute() error {
	cmd := NewCmdRoot()
	cmd.CompletionOptions.DisableDefaultCmd = true
	return cmd.Execute()
}
