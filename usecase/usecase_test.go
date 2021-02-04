package usecase

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Eleron8/filestackTestApi/getfile"
	"github.com/Eleron8/filestackTestApi/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestUsecase(t *testing.T) {
	imageUrl := "https://static.wikia.nocookie.net/zelda_gamepedia_en/images/3/35/WW_Link_3.png/revision/latest?cb=20130913013026"
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
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
	logger, _ := loggerCfg.Build()
	fHandl := getfile.NewFileHandler(httpClient, logger)
	folderName := "testImages"
	maxgoroutines := 5
	transforms := []models.Transform{
		{
			Type: models.Rotate,
			Params: models.Param{
				Degrees: 30.0,
				Width:   0,
				Height:  0,
			},
		},
		{
			Type: models.Crop,
			Params: models.Param{
				Degrees: 0,
				Width:   600,
				Height:  400,
			},
		},
		{
			Type: models.RemoveExif,
			Params: models.Param{
				Degrees: 0,
				Width:   0,
				Height:  0,
			},
		},
		{
			Type: models.Crop,
			Params: models.Param{
				Degrees: 0,
				Width:   200,
				Height:  100,
			},
		},
		{
			Type: models.Rotate,
			Params: models.Param{
				Degrees: 90,
				Width:   0,
				Height:  0,
			},
		},
		{
			Type: models.Crop,
			Params: models.Param{
				Degrees: 0,
				Width:   400,
				Height:  400,
			},
		},
	}
	dataTransform := models.TransformData{
		FileURL:    imageUrl,
		Transforms: transforms,
	}
	uCase := NewUsecase(fHandl, folderName, maxgoroutines, logger)
	f, err := os.Create("testarchive.zip")
	assert.Equal(t, nil, err, "create test archive error")
	err = uCase.FileFlow(dataTransform, f)
	assert.Equal(t, nil, err, "file flow error")
	fStat, err := os.Stat("testarchive.zip")
	assert.Equal(t, nil, err, "reading directory error")
	assert.NotEqual(t, 0, fStat.Size(), "file size")
}
