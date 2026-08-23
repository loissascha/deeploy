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

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the agent.",
	Long:  `Starts the agent and runs automatic jobs and connects to the control server defined.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// load default config path
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
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
	serveCmd.Flags().StringP("config", "c", "", "optional: path to agents config file")

	rootCmd.AddCommand(serveCmd)
}
