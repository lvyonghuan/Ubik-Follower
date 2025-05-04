package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

type UFollower struct {
	UUID string //Unique identifier for the follower

	Config Config
	Log    *ulog.Log
	OpType string

	status int //Ubik Follower's status

	r *gin.Engine

	plugin *plugin

	runtimeNodes map[int]*RuntimeNode
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

func (engine *UFollower) GetPluginList() map[string]*Plugin {
	return engine.plugin.plugins
}
