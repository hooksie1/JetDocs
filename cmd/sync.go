/*
Copyright © 2022 John Hooks

*/
package cmd

import (
	"fmt"
	"github.com/hooksie1/jetdocs/backend"
	"github.com/hooksie1/jetdocs/sync"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "syncs local markdown docs with jetstream",
	RunE:  syncData,
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolP("all", "a", false, "Sync all documents in directory")
	viper.BindPFlag("all", syncCmd.Flags().Lookup("all"))
	syncCmd.Flags().StringP("directory", "d", "./", "Directory to sync")
	viper.BindPFlag("directory", syncCmd.Flags().Lookup("directory"))
}

func syncData(cmd *cobra.Command, args []string) error {
	opts := getOptions()
	urls := viper.GetString("nats-urls")
	if urls == "" {
		urls = "nats://127.0.0.1:44566"
	}

	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		return err
	}

	b := backend.Nats{
		Conn: nc,
	}

	syncer := sync.FileSync{
		Directory: viper.GetString("directory"),
	}

	if viper.GetBool("all") {
		if err := syncer.ReadAllFiles(); err != nil {
			return err
		}

		if err := syncer.Sync(b); err != nil {
			return err
		}

		return nil

	}

	if args[0] == "" {
		return fmt.Errorf("file name must be supplied")
	}

	if err := syncer.ReadFile(args[0]); err != nil {
		return err
	}

	syncer.Sync(b)

	return nil

}
