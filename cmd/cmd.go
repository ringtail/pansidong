package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/ringtail/pansidong/types"
	"io/ioutil"
	"encoding/json"
	"github.com/ringtail/pansidong/server"
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
		log.Fatal(err)
	}
}

func Run(cmd *cobra.Command, args []string) {
	c := &types.Config{
		GlobalConfig: &types.GlobalConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
	if cfgFile == "" {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.Fatal("Failed to read configFile,because of %s", err.Error())
		}
		err = json.Unmarshal(data, c)
		if err != nil {
			log.Fatal("Failed to create config,because of %s", err.Error())
		}
	}
	s := server.NewServer(c)
	s.Run()
}
