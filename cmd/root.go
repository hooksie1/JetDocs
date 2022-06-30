/*
Copyright © 2022 John Hooks

*/
package cmd

import (
	"os"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "jetdocs",
	Short: "Sync markdown notes with Jetstream",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringP("nats-urls", "n", "", "The NATS urls")
	viper.BindPFlag("nats-urls", rootCmd.PersistentFlags().Lookup("nats-urls"))
	rootCmd.PersistentFlags().StringP("creds", "c", "", "Creds for NATS")
	viper.BindPFlag("creds", rootCmd.PersistentFlags().Lookup("creds"))
}

func getOptions() []nats.Option {
	var opts []nats.Option
	if viper.GetString("creds") != "" {
		opts = append(opts, nats.UserCredentials(viper.GetString("creds")))
	}

	return opts
}
