/*
Copyright © 2022 John Hooks
*/
package cmd

import (
	"log"

	"github.com/hooksie1/jetdocs/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	RunE:  start,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntP("port", "p", 8080, "Server port")
	viper.BindPFlag("port", startCmd.Flags().Lookup("port"))
	startCmd.Flags().StringP("store-dir", "s", "./jetdocs-data", "Directory for storage")
	viper.BindPFlag("store-dir", startCmd.Flags().Lookup("store-dir"))
}

func start(cmd *cobra.Command, args []string) error {
	opts, err := getOptions()
	if err != nil {
		return err
	}

	var nc *nats.Conn

	if viper.GetString("nats-urls") == "" {
		nc, err = server.StartEmbeddedNATS(nc, viper.GetString("store-dir"))
		if err != nil {
			return err
		}

		if err := server.InitializeBucket(nc); err != nil {
			return err
		}

	} else {
		nc, err = nats.Connect(viper.GetString("nats-urls"), opts...)
		if err != nil {
			return err
		}
	}

	s := server.NewServer().SetNatsConn(nc).SetPort(viper.GetInt("port"))
	log.Fatal(s.Serve())

	return nil

}
