package main

import (
	data "github.com/PicPay/ms-data-formatter/core/v1/data"
	"github.com/PicPay/ms-data-formatter/pkg/log"
	"github.com/PicPay/ms-data-formatter/pkg/newrelic"
	"github.com/PicPay/ms-data-formatter/pkg/server"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/spf13/cobra"
)

var config server.Config

func checkFatal(err error) {
	if err != nil {
		log.Fatal("Error in application", err, nil)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		println(err.Error())
	}

	err = envconfig.Process("data_formatter", &config)
	checkFatal(err)

	server, err := server.New(&config, "/health")
	checkFatal(err)

	server.NewRelic = newrelic.Setup(config.AppEnv, config.NewRelicLicenseKey)
	if server.NewRelic.App != nil {
		server.HttpServer.Router.Use(nrgin.Middleware(server.NewRelic.App))
	}

	log.Info("Loading handlers for", nil)
	err = server.Load(
		"/v1",
		&data.Handler{},
	)

	checkFatal(err)

	var runSeed bool

	rootCmd := &cobra.Command{
		Use:                   "data-formatter-service [-hs]",
		Short:                 "data-formatter-service",
		Version:               "0.0.2",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			log.Info("Starting HTTP serverApp", &log.LogContext{
				"address": server.Config.HttpAddress,
			})

			server.Start()
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&runSeed, "seed", "s", false, "run database seed")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error to execute root command", err, nil)
	}
}
