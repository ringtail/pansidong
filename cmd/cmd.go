package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pansidong",
	Short: "pansidong(盘丝洞）is a proxy ips' store service.",
	Run:   Run,
}

var (
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "default config file.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}

func Run(cmd *cobra.Command, args []string) {
	// Do Stuff Here

}
