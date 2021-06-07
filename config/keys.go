package config

const (
	AppName = "app.name"
	Secret  = "secret"

	LoggerCode         = "logger.code"
	LoggerLevel        = "logger.level"
	LoggerEnableCaller = "logger.enable.caller"

	LoggerZapEncoding                   = "logger.zap.encoding"
	LoggerZapDevelopment                = "logger.zap.development"
	LoggerZapEncoderConfigKeyMessage    = "logger.zap.encoder.config.key.message"
	LoggerZapEncoderConfigKeyLevel      = "logger.zap.encoder.config.key.level"
	LoggerZapEncoderConfigKeyTime       = "logger.zap.encoder.config.key.time"
	LoggerZapEncoderConfigKeyName       = "logger.zap.encoder.config.key.name"
	LoggerZapEncoderConfigKeyCaller     = "logger.zap.encoder.config.key.caller"
	LoggerZapEncoderConfigKeyStacktrace = "logger.zap.encoder.config.key.stacktrace"

	MongoDBUrl         = "mongo.db.url"
	MongoDBName        = "mongo.db.name"
	MongoDBTimeOut     = "mongo.db.timeout"
	MongoDBMinPoolSize = "mongo.db.pool.size.min"
	MongoDBMaxPoolSize = "mongo.db.pool.size.max"

	AcelleMailURI   = "acelle.mail.uri"
	AcelleMailToken = "acelle.mail.token"
)
