package factory

import (
	"github.com/Cliengo/acelle-mail/config"
	zap2 "github.com/Cliengo/acelle-mail/container/logger/factory/zap"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var logFactoryBuilderMap = map[string]logFactoryInterface{
	"zap": &zap2.ZapFactory{},
}

type logFactoryInterface interface {
	Build() error
}

func GetLogFactoryBuilder(key string) logFactoryInterface {
	return logFactoryBuilderMap[key]
}

func init() {
	code := viper.GetString(config.LoggerCode)
	if err := GetLogFactoryBuilder(code).Build(); err != nil {
		errors.Wrap(err, "fail to load logger")
	}
}
