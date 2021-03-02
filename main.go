package main

import (
	"github.com/PicPay/ms-data-crawler/core/v1/configuration"
	"github.com/PicPay/ms-data-crawler/pkg/log"
	"github.com/PicPay/ms-data-crawler/pkg/newrelic"
	"github.com/PicPay/ms-data-crawler/pkg/server"

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
	err := godotenv.Load(".env", ".env.testing")
	if err != nil {
		println(err.Error())
	}

	err = envconfig.Process("dc", &config)

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
		&configuration.Handler{},
	)

	checkFatal(err)

	rootCmd := &cobra.Command{
		Use:                   "data-crawler-service [-hs]",
		Short:                 "data-crawler-service",
		Version:               "0.0.1",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			log.Info("Starting HTTP serverApp", &log.LogContext{
				"address": server.Config.HttpAddress,
			})

			server.Start()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error to execute root command", err, nil)
	}
}
