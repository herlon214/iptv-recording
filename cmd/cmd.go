package cmd

import (
	"github.com/herlon214/iptv-recording/cmd/kubernetes"
	"github.com/herlon214/iptv-recording/cmd/local"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "iptv-rec",
	Short: "Record IPTV using ffmpeg",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(local.LocalCmd)
	RootCmd.AddCommand(kubernetes.KubernetesCmd)
}
