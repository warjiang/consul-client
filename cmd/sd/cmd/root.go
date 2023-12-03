package cmd

import "github.com/spf13/cobra"

var (
	ConsulHost string
)

func NewServiceDiscoveryCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "sd",
		Short: "sd tools for consul",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	cmds.AddCommand(NewCmdLookUp())
	return cmds
}
