package engine

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lvyonghuan/Ubik-Util/uerr"
)

//Run the plugin which is mounted
//And also stop plugins

// startPlugin starts the plugin
// It will start a bat or a shell script which same as the plugin's name
func startPlugin(plugin *Plugin, e *UFollower) error {
	//get the op sys type
	opType := e.OpType

	switch opType {
	case Windows:
		return startPluginOnWindows(plugin, e)

	case Linux:
		return startPluginOnLinux(plugin, e)

	case Mac:
		return startPluginOnMac(plugin, e)

	default:
		return uerr.NewError(errors.New("unsupported operating system"))
	}
}

func startPluginOnWindows(plugin *Plugin, engine *UFollower) error {
	pluginName, port, pluginPath := getStartInfo(plugin, engine)

	//stitch the path
	pluginScriptPath := splicePath(pluginPath, pluginName, winSuffix)

	//execute the bat file
	//passing in port parameters
	absPath, err := filepath.Abs(pluginScriptPath)
	if err != nil {
		return uerr.NewError(err)
	}
	cmd := exec.Command(absPath, port)
	cmd.Dir = spliceWorkDir(pluginPath, pluginName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func startPluginOnLinux(plugin *Plugin, engine *UFollower) error {
	pluginName, port, pluginPath := getStartInfo(plugin, engine)

	//stitch the path
	pluginPath = splicePath(pluginPath, pluginName, linuxSuffix)

	//execute the shell file
	//passing in port parameters
	cmd := exec.Command(pluginPath, port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

// FIXME I don't have a mac, so I can't test it
func startPluginOnMac(plugin *Plugin, engine *UFollower) error {
	pluginName, port, pluginPath := getStartInfo(plugin, engine)

	//stitch the path
	pluginPath = splicePath(pluginPath, pluginName, macSuffix)

	//execute the shell file
	//passing in port parameters
	cmd := exec.Command(pluginPath, port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func getStartInfo(plugin *Plugin, engine *UFollower) (string, string, string) {
	//get the follower running port and plugin save path
	port, pluginPath := engine.Config.Port, engine.Config.PluginPath

	//get the plugin name
	pluginName := plugin.PluginMetaData.Name

	return pluginName, port, pluginPath
}

const (
	winSuffix   = ".bat"
	linuxSuffix = ".sh"
	macSuffix   = ".sh"
)

func spliceWorkDir(pluginPath, pluginName string) string {
	return pluginPath + pluginName + "/"
}

func splicePath(pluginPath, pluginName, suffix string) string {
	return spliceWorkDir(pluginPath, pluginName) + pluginName + suffix
}
