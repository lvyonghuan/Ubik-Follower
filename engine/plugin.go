package engine

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

type Plugin struct {
	PluginMetaData *uplugin.Plugin

	runtimeNodes map[int]*RuntimeNode
	//The number of nodes that are mounted
	MountNodeNum int
	//when a mounted plugin is running, will send a signal to this channel to break it's blocking
	WaitRunningBlockingChannel chan struct{}

	checkMutex sync.Mutex
}

// RuntimeNode
// When a node is mounted, it will be converted to a runtimeNode
// A node can be mounted multiple times, and each mount will generate a runtimeNode
// Runtime node id distributed by leader
type RuntimeNode struct {
	ID       int    `json:"id"`
	NodeName string `json:"node_name"`

	pluginInfo *Plugin

	params      uplugin.Params
	OutputEdges map[string][]Edge `json:"output_edges"` //A output port can have multiple output edges
}

// Edge is a directed edge, from producer to consumer
type Edge struct {
	//The addr point to the consumer PLUGIN.
	//If addr is not nil, node can send output directly to the next node
	Addr string `json:"addr"`

	ProducerID       int    `json:"producer_id"`
	ConsumerID       int    `json:"consumer_id"`
	ProducerPortName string `json:"producer_port_name"`
	ConsumerPortName string `json:"consumer_port_name"`

	channel chan []byte
}

// Check if the plugin is already mounted
// If not, then will mount the plugin
func (engine *UFollower) checkMount(plugin *Plugin) {
	plugin.checkMutex.Lock()
	if _, isExist := engine.plugin.mountedPlugins[plugin.PluginMetaData.Name]; isExist {
		engine.Log.Debug(fmt.Sprintf("plugin %s already mounted", plugin.PluginMetaData.Name))

		plugin.MountNodeNum++
		plugin.checkMutex.Unlock()
		return
	} else {
		engine.Log.Debug(fmt.Sprintf("mounting plugin %s", plugin.PluginMetaData.Name))

		engine.plugin.mountedPlugins[plugin.PluginMetaData.Name] = plugin
		plugin.MountNodeNum++
		plugin.checkMutex.Unlock()

		//add a block for prepare func
		//when plugin run and init, will send a signal to this channel to break the block
		plugin.WaitRunningBlockingChannel = make(chan struct{}, 1)
		engine.plugin.waitRunningBlockingChannels = append(engine.plugin.waitRunningBlockingChannels, plugin.WaitRunningBlockingChannel)

		//start the plugin
		err := startPlugin(plugin, engine)
		if err != nil {
			engine.Log.Error(err)
			//TODO 传出错误
		}

	}
}

// Check if the plugin can be unmounted
func (engine *UFollower) checkUnmount(plugin *Plugin) {
	plugin.checkMutex.Lock()
	plugin.MountNodeNum--
	if plugin.MountNodeNum <= 0 {
		//TODO unmount plugin
		//...

		delete(engine.plugin.mountedPlugins, plugin.PluginMetaData.Name)
	}

	plugin.checkMutex.Unlock()
}

//// When edge wants follower to forward its message,
//// a disguised runtime node will be generated
//func (engine *UFollower) generateDisguisedRuntimeNode() *RuntimeNode {
//
//}

const infoFileName = "info.json"

func (engine *UFollower) detectPlugins() error {
	pluginPath := engine.Config.PluginPath
	plugins := make(map[string]*Plugin)
	//first, scan the plugin path, find all dir
	dir, err := os.Open(pluginPath)
	if err != nil {
		return uerr.NewError(err)
	}

	//scan to find all dir
	dirs, err := dir.Readdir(0)
	if err != nil {
		return uerr.NewError(err)
	}

	//then open each dir, read the info.json file
	for _, d := range dirs {
		if d.IsDir() {
			//stitch paths
			infoPath := pluginPath + "/" + d.Name() + "/" + infoFileName

			//parse info.json
			pluginMetaData, err := uplugin.ReadPluginInfo(infoPath)
			if err != nil {
				engine.Log.Error(errors.New(fmt.Sprintf("read plugin info.json error: %s", err.Error())))
				continue
			}

			//create plugin
			plugin := Plugin{
				PluginMetaData: pluginMetaData,
				runtimeNodes:   make(map[int]*RuntimeNode),
			}

			//add to plugins
			plugins[plugin.PluginMetaData.Name] = &plugin

			engine.Log.Info(fmt.Sprintf("plugin %s detected", plugin.PluginMetaData.Name))
		}
	}
	engine.plugin.plugins = plugins

	return nil
}

func (engine *UFollower) PutParams(id int, paramsJson []byte) error {
	var params uplugin.Params
	err := ujson.Unmarshal(paramsJson, &params)
	if err != nil {
		return uerr.NewError(err)
	}

	//Check if the id already exists
	if runtimeNode, isExist := engine.runtimeNodes[id]; !isExist {
		return uerr.NewError(errors.New("runtime node" + fmt.Sprintf("%d", id) + "not exist"))
	} else {
		runtimeNode.params = params
	}

	return nil
}

func (engine *UFollower) InitPluginsNodes() {
	for _, plugin := range engine.plugin.mountedPlugins {
		engine.initPluginNodes(plugin)
	}
}

// WaitPrepare will block until all plugins are prepared
func (engine *UFollower) WaitPrepare() {
	engine.Log.Info("Start waiting for plugins to prepare...")

	blockChannels := engine.plugin.waitRunningBlockingChannels
	for _, blockChan := range blockChannels {
		//wait for the signal
		<-blockChan
	}

	engine.Log.Info("All plugins are prepared, continue to run plugins...")
}

func (engine *UFollower) RunPlugins() error {
	engine.Log.Info("Start running plugins...")
	for _, plugin := range engine.plugin.mountedPlugins {
		err := engine.callPluginToRun(plugin)
		if err != nil {
			engine.Log.Error(err)
			continue
		}
	}
	return nil
}
