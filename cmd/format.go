package cmd

import (
	"os"

	"github.com/pszeto/access-log-formatter/pkg/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfg = &format.Config{}

var rootCmd = &cobra.Command{
	Use:   "format-access-log",
	Short: "Formats default Envoy / Istio Access Logs",
	Long:  "Formats default Envoy / Istio Access Logs",
	Run: func(cmd *cobra.Command, args []string) {
		format.New(cfg).Entry()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfg.File, "file", "", "Specify access log file. default to prompt if not set")
}
