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
	rootCmd.PersistentFlags().String("nkey", "", "NKey Seed File")
	viper.BindPFlag("nkey", rootCmd.PersistentFlags().Lookup("nkey"))
	rootCmd.PersistentFlags().String("tlscert", "", "TLS client certificate file")
	viper.BindPFlag("tlscert", rootCmd.PersistentFlags().Lookup("tlscert"))
	rootCmd.PersistentFlags().String("tlskey", "", "Private key file for client certificate")
	viper.BindPFlag("tlskey", rootCmd.PersistentFlags().Lookup("tlskey"))
	rootCmd.PersistentFlags().String("tlscacert", "", "CA certificate to verify peer against")
	viper.BindPFlag("tlscacert", rootCmd.PersistentFlags().Lookup("tlscacert"))
}

func getOptions() ([]nats.Option, error) {
	var opts []nats.Option
	if viper.GetString("creds") != "" {
		opts = append(opts, nats.UserCredentials(viper.GetString("creds")))
	}

	cCert, cKey := viper.GetString("tlscert"), viper.GetString("tlskey")
	if cCert != "" && cKey != "" {
		opts = append(opts, nats.ClientCert(cCert, cKey))
	}

	if viper.GetString("tlscacert") != "" {
		opts = append(opts, nats.RootCAs(viper.GetString("tlscacert")))
	}

	if viper.GetString("nkey") != "" {
		opt, err := nats.NkeyOptionFromSeed(viper.GetString("nkey"))
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}
