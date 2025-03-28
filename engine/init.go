package engine

import (
	"errors"
	"os"
	"runtime"

	"github.com/lvyonghuan/Ubik-Util/uconfig"
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

func InitEngine(inTest bool) *UFollower {
	var config Config
	var err error

	if !inTest {
		config, err = readConfig(configPath)
		if err != nil {
			fatalErrHandel(err)
		}
	} else {
		config, err = readConfig(testConfigPath)
		if err != nil {
			fatalErrHandel(err)
		}
	}

	engine, err := makeUFollower(config)
	if err != nil {
		fatalErrHandel(err)
	}

	engine.Log.Info("Scanning plugins...")
	err = engine.detectPlugins() //detect installed plugins
	if err != nil {
		engine.Log.Fatal(err)
		os.Exit(1)
	}

	engine.detectOpType() //detect the operating system type
	engine.Log.Debug("Operating system type: " + engine.OpType)

	return engine
}

func readConfig(configPath string) (Config, error) {
	var config Config
	err := uconfig.Read(configPath, &config)
	if err != nil {
		return config, uerr.NewError(err)
	}
	return config, nil
}

func initLog(log Log) *ulog.Log {
	l := ulog.Log{
		Level:       log.Level,
		IsSave:      log.IsSave,
		LogSavePath: log.LogSavePath,
	}
	l.InitLog()

	return &l
}

func (engine *UFollower) detectOpType() {
	// Detect the operating system type
	if runtime.GOOS == Windows {
		engine.OpType = Windows
	} else if runtime.GOOS == Linux {
		engine.OpType = Linux
	} else if runtime.GOOS == Mac {
		engine.OpType = Mac
	} else {
		fatalErrHandel(uerr.NewError(errors.New("unsupported operating system " + runtime.GOOS)))
	}
}

func makeUFollower(config Config) (*UFollower, error) {
	log := initLog(config.Log)
	return &UFollower{
		Config:       config,
		Log:          log,
		runtimeNodes: make(map[int]*RuntimeNode),
		plugin: &plugin{
			plugins:        make(map[string]*Plugin),
			mountedPlugins: make(map[string]*Plugin),
		},
	}, nil
}

func fatalErrHandel(err error) {
	l := initLog(Log{
		Level:       ulog.Debug,
		IsSave:      true,
		LogSavePath: "./logs",
	})
	l.Fatal(err)
	os.Exit(1)
}
