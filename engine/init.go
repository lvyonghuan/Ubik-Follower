package engine

import (
	"errors"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/lvyonghuan/Ubik-Util/uconfig"
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

func InitEngine(inTest bool) *UFollower {
	var config Config
	var err error

	if !inTest {
		config, err = readConfig(confPath + "/config.json")
		if err != nil {
			fatalErrHandel(err)
		}
	} else {
		config, err = readConfig(testConfPath + "/config.json")
		if err != nil {
			fatalErrHandel(err)
		}
	}

	engine, err := makeUFollower(config, inTest)
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

	//Connect to the leader
	engine.detectLeader()
	//Post the plugin list to leader
	err = engine.postPlugins()
	if err != nil { //TODO：重试机制
		engine.Log.Fatal(err)
		os.Exit(1)
	}
	//Start heartbeat
	err = engine.initHeartbeat()
	if err != nil {
		engine.Log.Fatal(err)
		os.Exit(1)
	}

	return engine
}

func readConfig(configPath string) (Config, error) {
	var config Config
	err := uconfig.Read(configPath, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func initLog(config Config, uuid string) ulog.Log {
	l := ulog.NewLogWithPost(config.Log.Level, config.Log.IsSave, config.Log.LogSavePath, config.LeaderAddr, uuid)

	return l
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

func (engine *UFollower) getUUID(uuidPath string) error {
	//Check if the file exists
	_, err := os.Stat(uuidPath)
	if err != nil {
		if os.IsNotExist(err) {
			//Create UUID for follower
			err = engine.createUUID(uuidPath)
			if err != nil {
				return err
			}

			return nil
		} else {
			return uerr.NewError(err)
		}
	}

	content, err := os.ReadFile(uuidPath)
	if err != nil {
		return uerr.NewError(err)
	}

	engine.UUID = string(content)
	return nil
}

func (engine *UFollower) createUUID(uuidPath string) error {
	//Create the file
	file, err := os.Create(uuidPath)
	if err != nil {
		return uerr.NewError(err)
	}
	defer file.Close()

	//Generate a new UUID
	u := uuid.New()

	//Write the UUID to the file
	_, err = file.WriteString(u.String())
	if err != nil {
		return uerr.NewError(err)
	}

	engine.UUID = u.String()
	return nil
}

func makeUFollower(config Config, inTest bool) (*UFollower, error) {
	//Initialize the UFollower engine
	engine := &UFollower{
		Config:    config,
		workflows: make(map[string]runtimeNodes),
		plugin: &plugin{
			plugins:        make(map[string]*Plugin),
			mountedPlugins: make(map[string]*Plugin),
		},
	}

	//Initialize the UUID
	if !inTest {
		err := engine.getUUID(confPath + "/uuid")
		if err != nil {
			fatalErrHandel(err)
			os.Exit(1)
		}
	} else {
		err := engine.getUUID(testConfPath + "/uuid")
		if err != nil {
			fatalErrHandel(err)
			os.Exit(1)
		}
	}

	//Initialize the log
	log := initLog(config, engine.UUID)
	engine.Log = log

	return engine, nil
}

// handel error when log isn't initialized
func fatalErrHandel(err error) {
	l := ulog.NewLogWithoutPost(ulog.Debug, true, "./logs")
	l.Fatal(err)
	os.Exit(1)
}
