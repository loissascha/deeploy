/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"local/deeploy/internal/agent"
	"local/deeploy/internal/settings"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
)

var configPath string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// load default config path
		if configPath == "" {
			var err error
			configPath, err = settings.GetAgentPath()
			if err != nil {
				panic(err)
			}
		}

		// load agent settings
		agentSettings, err := settings.LoadAgentSettings(configPath)
		if err != nil {
			panic(err)
		}
		slog.Info("agent settings:", "as", agentSettings)

		// create agent
		a := agent.NewAgent(agentSettings)

		// run agent connection to server
		go func() {
			for {
				if err := a.RunConn(ctx); err != nil {
					slog.Error("error running the server conn", "err", err)
				}
				slog.Info("retrying in 10 seconds...")
				time.Sleep(10 * time.Second)
			}
		}()

		// TODO: load jobs and run jobs when triggered

		for {
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "", "optional: path to agents config file")

	rootCmd.AddCommand(serveCmd)
}
