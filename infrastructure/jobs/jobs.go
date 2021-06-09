package jobs

import "github.com/Cliengo/acelle-mail/container/logger"

type Job struct {
	logger logger.Logger
}

func New(logger logger.Logger) Job {
	return Job{
		logger: logger,
	}
}

func (jb Job) Run() {
	jb.logger.Info("Hello world")
}
