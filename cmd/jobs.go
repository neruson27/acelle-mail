package cmd

import (
	"github.com/Cliengo/acelle-mail/container"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/spf13/cobra"
)

var jobsCmd = &cobra.Command{
	Use: "jobs",
	Run: runJobCmd,
}

func runJobCmd(cmd *cobra.Command, args []string) {
	logger.Log.Info("Starting jobs command")
	jobs := container.NewJobs()
	jobs.Run()
}

func init() {
	rootCmd.AddCommand(jobsCmd)
}
