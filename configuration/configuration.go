package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Config Cfg
)

type Cfg struct {
	ServerPort    string
	MaxGoroutines int
	FolderName    string
	ProjectID     string
	BucketName    string
}

func init() {
	loggerCfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	Logger, _ = loggerCfg.Build()
	viper.SetConfigName("config")
	viper.AddConfigPath("./configuration")
	err := viper.ReadInConfig()
	if err != nil {
		Logger.Fatal("can't read config file", zap.Error(err))
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		Logger.Fatal("unable to decode config in struct", zap.Error(err))
	}

}
