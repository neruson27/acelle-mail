package zap

import (
	"github.com/Cliengo/acelle-mail/config"
	"github.com/Cliengo/acelle-mail/container/logger"
	"github.com/pkg/errors"
	conf "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapFactory struct{}

var ZapLogger *zap.Logger

func (zf ZapFactory) Build() error {
	zapConfig, err := zf.getConfig()
	if err != nil {
		return err
	}
	log, err := zapConfig.Build()
	if err != nil {
		return errors.Wrap(err, "fail to build logger")
	}
	log.Debug("Logger construction succeeded")
	defer log.Sync()
	zSugarLog := log.Sugar()
	ZapLogger = log
	logger.SetLogger(zSugarLog)
	return nil
}

func (zf ZapFactory) getConfig() (zap.Config, error) {
	level := zap.NewAtomicLevel().Level()
	if err := level.Set(conf.GetString(config.LoggerLevel)); err != nil {
		return zap.Config{}, errors.Wrap(err, "fail to load level type")
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: conf.GetBool(config.LoggerZapDevelopment),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         conf.GetString(config.LoggerZapEncoding),
		EncoderConfig:    zf.newEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}, nil
}

func (zf ZapFactory) newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     conf.GetString(config.LoggerZapEncoderConfigKeyMessage),
		LevelKey:       conf.GetString(config.LoggerZapEncoderConfigKeyLevel),
		TimeKey:        conf.GetString(config.LoggerZapEncoderConfigKeyTime),
		NameKey:        conf.GetString(config.LoggerZapEncoderConfigKeyName),
		CallerKey:      conf.GetString(config.LoggerZapEncoderConfigKeyCaller),
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  conf.GetString(config.LoggerZapEncoderConfigKeyStacktrace),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
