package main

import (
	"github.com/Cliengo/acelle-mail/cmd"
	"github.com/Cliengo/acelle-mail/container/logger"
	_ "github.com/Cliengo/acelle-mail/container/logger/factory"
	"os"
)

func main() {
	if logger.Log == nil {
		os.Exit(1)
	}
	cmd.Execute()
}
