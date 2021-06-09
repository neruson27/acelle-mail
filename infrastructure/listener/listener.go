package listener

import (
	"context"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/repository"
	"github.com/Cliengo/acelle-mail/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type (
	streamProcessor struct {
		Key      string
		Watcher  func(resumeToken *primitive.ObjectID) (*mongo.ChangeStream, error)
		Executor func(stream bson.M) error
	}

	StreamsListener struct {
		logger     logger.Logger
		repository repository.StreamWatcherRepository
		service    services.StreamsService
	}
)

func New(logger logger.Logger, repository repository.StreamWatcherRepository, service services.StreamsService) StreamsListener {
	return StreamsListener{
		logger:     logger,
		repository: repository,
		service:    service,
	}
}

func (sl StreamsListener) Run(ctx context.Context) {
	streamsProcessors := sl.retrieveStreamsProcessors()
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(streamsProcessors))
	for _, streamProcessor := range streamsProcessors {
		go sl.execute(ctx, waitGroup, streamProcessor)
	}
	waitGroup.Wait()
}

// Retorna el listado de los watchers para los stream changes
func (sl StreamsListener) retrieveStreamsProcessors() []streamProcessor {
	return []streamProcessor{
		//{Key: "company-update", Watcher: sl.repository.WatcherUpdatesCompanyPlan, Executor: sl.service.ProcessUpdateCompany},
		{Key: "contacts", Watcher: sl.repository.WatcherContactEvents, Executor: sl.service.ProcessContactEvent},
	}
}

func (sl StreamsListener) execute(ctx context.Context, waitGroup sync.WaitGroup, processor streamProcessor) {
	defer waitGroup.Done()
	sl.logger.Infof("Execute, key: %s", processor.Key)
	resumeToken, err := sl.repository.RetrieveLastResumeToken(processor.Key)
	if err != nil {
		sl.logger.Errorf("%s, fail to retrieve last resumeToken, error: %s", processor.Key, err)
		resumeToken = nil
	}

	sl.logger.Infof("Retrieve cs, key: %s", processor.Key)
	cs, err := processor.Watcher(resumeToken)
	if err != nil {
		sl.logger.Errorf("%s, fail to connect to change streams, error: %s", processor.Key, err)
		//TODO: Ver que se puede hacer en este punto si no se puede conectar al change stream
		return
	}

	defer cs.Close(ctx)
	sl.logger.Infof("Starting cs, key: %s", processor.Key)
	for cs.Next(ctx) {
		sl.logger.Info(cs.ResumeToken())
		var data bson.M
		if err = cs.Decode(&data); err != nil {
			sl.logger.Errorf("fail to process stream, %s", err)

		} else {
			//TODO: Pendiente de manejar el Resumen token y el err que se genera
			_ = processor.Executor(data)
		}
	}

}
