/*
Copyright © 2022 John Hooks

*/
package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hooksie1/jetdocs/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
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
	opts := getOptions()
	var nc *nats.Conn
	var err error

	if viper.GetString("nats-urls") == "" {
		nc, err = server.StartEmbeddedNATS(nc)
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

	r := mux.NewRouter()
	s := server.NewServer().SetNatsConn(nc).SetRouter(r)

	r.HandleFunc("/pages/{id}", s.GetPage).Methods("GET")
	r.HandleFunc("/pages", s.GetPages).Methods("GET")

	port := fmt.Sprintf(":%d", viper.GetInt("port"))
	s.Port = viper.GetInt("port")

	log.Fatal(http.ListenAndServe(port, r))

	return nil

}
