package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func pluginInit(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	pluginName := c.Query("plugin_name")
	if pluginName == "" {
		engine.Log.Error(errors.New("plugin_name is empty"))
		errorResponse(c, 400, "plugin_name is empty")
		return
	}

	pluginURL := c.Query("plugin_addr")
	if pluginURL == "" {
		engine.Log.Error(errors.New("plugin " + pluginName + "'s plugin_addr is empty"))
		errorResponse(c, 400, "plugin_addr is empty")
		return
	}

	plugin, ok := engine.GetPluginList()[pluginName]
	if !ok {
		engine.Log.Error(errors.New("plugin " + pluginName + " not found"))
		errorResponse(c, 400, "plugin not found")
		return
	}
	plugin.Addr = pluginURL

	//FIXME Shouldn't be here. It should be in pre.

	successResponse(c, "")
}
