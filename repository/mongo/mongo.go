package mongo

import (
	"context"
	"github.com/Cliengo/acelle-mail/config"
	"github.com/pkg/errors"
	conf "github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func New() (Client, error) {
	client, database, err := newClient()
	if err != nil {
		return Client{}, errors.Wrap(err, "")
	}

	return Client{
		client:   client,
		database: database,
	}, nil
}

func newClient() (*mongo.Client, *mongo.Database, error) {
	mongoDBURL, err := retrieveMongoURL()
	if err != nil {
		return nil, nil, err
	}

	cs, err := connstring.ParseAndValidate(mongoDBURL)
	if err != nil {
		return nil, nil, err
	}

	minPoolSize, maxPoolSize := retrieveMinMaxPoolSize()
	clOptions := options.Client()

	clOptions.ApplyURI(mongoDBURL)
	clOptions.SetMinPoolSize(minPoolSize)
	clOptions.SetMaxPoolSize(maxPoolSize)
	clOptions.SetAppName(conf.GetString(config.AppName))

	client, err := mongo.NewClient(clOptions)
	if err != nil {
		return nil, nil, errors.Wrap(err, "1")
	}

	ctx, cancel := context.WithTimeout(context.Background(), retrieveTimeOut())
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "2")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, nil, errors.Wrap(err, "3")
	}

	dataBaseName, err := retrieveDatabaseName(cs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "4")
	}
	database := client.Database(dataBaseName)
	return client, database, nil
}

func retrieveMongoURL() (string, error) {
	mongoURL := conf.GetString(config.MongoDBUrl)
	if mongoURL == "" {
		return "", errors.New("not valid mongo url")
	}
	return mongoURL, nil
}

func retrieveTimeOut() time.Duration {
	timeOut := conf.GetDuration(config.MongoDBTimeOut)

	if timeOut <= 0 {
		timeOut = 15
	}
	return timeOut * time.Second
}

func retrieveMinMaxPoolSize() (uint64, uint64) {
	minPoolSize := conf.GetUint64(config.MongoDBMinPoolSize)
	if minPoolSize <= 0 {
		minPoolSize = 5
	}

	maxPoolSize := conf.GetUint64(config.MongoDBMaxPoolSize)
	if maxPoolSize <= 0 {
		maxPoolSize = 15
	}
	return minPoolSize, maxPoolSize
}

func retrieveDatabaseName(cs connstring.ConnString) (string, error) {
	if cs.Database != "" {
		return cs.Database, nil
	}
	database := conf.GetString(config.MongoDBName)
	if database == "" {
		return "", errors.New("not database set to connect")
	}

	//TODO: Seria bueno estar verificando si la database existe en mongo
	return database, nil
}
