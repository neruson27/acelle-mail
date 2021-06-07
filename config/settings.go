package config

import (
	conf "github.com/spf13/viper"
	"time"
)

func init() {
	setSettings()
}

func setSettings() {
	defaultAppConfig := map[string]interface{}{
		AppName: "cliengo-acelle-mail",
		Secret:  "your-256-bit-secret",

		LoggerCode:         "zap",
		LoggerLevel:        "debug",
		LoggerEnableCaller: true,

		LoggerZapEncoding:                "console",
		LoggerZapDevelopment:             true,
		LoggerZapEncoderConfigKeyMessage: "msg",
		LoggerZapEncoderConfigKeyLevel:   "level",
		LoggerZapEncoderConfigKeyTime:    "ts",
		LoggerZapEncoderConfigKeyName:    "logger",
		LoggerZapEncoderConfigKeyCaller:  "fn",

		//Prod
		//MongoDBUrl: "mongodb+srv://CLIENGOPROD_READONLY:AwUzKcGk3LO5Zv9Q@prod-cliengo-core.9szjk.mongodb.net/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&ssl=true",

		//Stage
		MongoDBUrl:         "mongodb+srv://CLIENGOSTAGE:9G9dPkHDBnnzasey@stage-paid.qdpwm.mongodb.net/convergency_prod_paid?authSource=admin&replicaSet=atlas-ny8za9-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true",
		MongoDBName:        "convergency_prod_paid",
		MongoDBTimeOut:     time.Duration(15),
		MongoDBMinPoolSize: uint64(5),
		MongoDBMaxPoolSize: uint64(15),

		AcelleMailURI:   "https://emailmkt.stagecliengo.com",
		AcelleMailToken: "k4FvkXNcIFPnDlYKiUO2VUkQYdwqIVot1iOxRC1K13sk0uTSWHWoLLXUDRKR",
	}

	for key, value := range defaultAppConfig {
		conf.SetDefault(key, value)
	}
}
