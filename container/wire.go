// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package container

import (
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/Cliengo/acelle-mail/infrastructure/http"
	"github.com/Cliengo/acelle-mail/infrastructure/listener"
	"github.com/Cliengo/acelle-mail/repository"
	"github.com/Cliengo/acelle-mail/repository/mongo"
	"github.com/Cliengo/acelle-mail/services"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	mongo.New,
	mongo.NewMongoIntegrationRepository,
	wire.Bind(new(repository.IntegrationRepository), new(mongo.IntegrationRepository)),
	mongo.NewMongoWatcherRepository,
	wire.Bind(new(repository.StreamWatcherRepository), new(mongo.WatcherRepository)),
)

var servicesSet = wire.NewSet(
	services.NewAcelleMailService,
	services.NewAccountService,
	services.NewStreamService,
)

var handlerSet = wire.NewSet(
	http.NewHealthHandler,
	http.NewIntegrationHandler,
	http.NewServerHandlers,
)

func NewServer() (http.Server, error) {
	wire.Build(logger.GetLogger, repositorySet, servicesSet, handlerSet, http.New)
	return http.Server{}, nil
}

func NewListener() (listener.StreamsListener, error) {
	wire.Build(logger.GetLogger, repositorySet, servicesSet, listener.New)
	return listener.StreamsListener{}, nil
}
