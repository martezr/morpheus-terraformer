package cmd

import (
	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/martezr/morpheus-terraformer/resources"
	"github.com/spf13/cobra"
)

func init() {
}

var installProviderCmd = &cobra.Command{
	Use:   "install",
	Short: "Print the version number of Morpheus Terraformer",
	Long:  `Show this help output, or the help for a specified subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := morpheus.NewClient("")
		client.SetAccessToken("", "", 0, "write")
		resources.InstallProvider()
	},
}
