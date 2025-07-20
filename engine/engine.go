package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/ulog"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

type UFollower struct {
	UUID string //Unique identifier for the follower

	Config Config
	Log    ulog.Log
	OpType string

	status int //Ubik Follower's status

	r *gin.Engine

	plugin *plugin

	runtimeNodes map[int]*RuntimeNode

	heartbeat *heartbeat
}

type plugin struct {
	plugins        map[string]*Plugin
	mountedPlugins map[string]*Plugin

	waitRunningBlockingChannels []chan struct{}
}

const (
	Windows = "windows"
	Linux   = "linux"
	Mac     = "darwin"
)

func (engine *UFollower) SetGinEngine(r *gin.Engine) {
	engine.r = r
}

func (engine *UFollower) GetPluginList() map[string]uplugin.Plugin {
	originPlugins := engine.plugin.plugins
	plugins := make(map[string]uplugin.Plugin, len(originPlugins))

	for name, p := range originPlugins {
		metaData := p.PluginMetaData

		plugins[name] = *metaData
	}

	return plugins
}
