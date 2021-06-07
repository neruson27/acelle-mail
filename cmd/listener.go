package cmd

import (
	"context"
	"github.com/Cliengo/acelle-mail/container"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var listenerCmd = &cobra.Command{
	Use: "listener",
	Run: runListenerCmd,
}

func runListenerCmd(cmd *cobra.Command, args []string) {
	listener, err := container.NewListener()
	if err != nil {
		logger.Log.Errorf("fail to initialize listener, error: %s", err)
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	go listener.Run(ctx)
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	logger.Log.Info("Shutting listener server...")
	cancelFunc()
	logger.Log.Info("Listener stopped")
}

func init() {
	rootCmd.AddCommand(listenerCmd)
}
