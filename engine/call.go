package engine

import (
	"bytes"
	"net/http"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
)

const (
	initNodes = "/init"
	setParams = "/param"
	run       = "/run"
)

// InitPluginNodes After plugin init, follower should ask plugin to add nodes and edges
// TODO: Error handel
func (engine *UFollower) initPluginNodes(plugin *Plugin) {
	nodes := plugin.runtimeNodes

	//Marshal the nodes data
	jsonData, err := ujson.Marshal(nodes)
	if err != nil {
		engine.Log.Error(uerr.NewError(err))
		return
	}

	url := plugin.PluginMetaData.Uri + initNodes
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		engine.Log.Error(uerr.NewError(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	//TODO: Check the response
	_, err = engine.callPlugin(req)
	if err != nil {
		engine.Log.Error(uerr.NewError(err))
		return
	}

	//Set Params of the nodes
	err = engine.setParamsToRuntimeNodes(nodes)
	if err != nil {
		engine.Log.Error(uerr.NewError(err))
		return
	}

	plugin.WaitRunningBlockingChannel <- struct{}{} //break the blocking

	engine.Log.Debug("Init plugin " + plugin.PluginMetaData.Name + " nodes success")
}

func (engine *UFollower) setParamsToRuntimeNodes(nodes map[int]*RuntimeNode) error {
	for _, node := range nodes {
		//If >0, prof it has params, then set the params
		if len(node.params) > 0 {
			url := node.pluginInfo.PluginMetaData.Uri + setParams
			jsonData, err := ujson.Marshal(node.params)
			if err != nil {
				engine.Log.Error(uerr.NewError(err))
				continue
			}

			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
			if err != nil {
				engine.Log.Error(uerr.NewError(err))
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			//TODO: Check the response
			_, err = engine.callPlugin(req)
			if err != nil {
				engine.Log.Error(uerr.NewError(err))
				continue
			}
		}
	}
	return nil
}

func (engine *UFollower) callPluginToRun(plugin *Plugin) error {
	url := plugin.PluginMetaData.Uri + run
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return uerr.NewError(err)
	}
	//TODO Check the response
	_, err = engine.callPlugin(req)
	if err != nil {
		return uerr.NewError(err)
	}
	return nil
}

func (engine *UFollower) callPlugin(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, uerr.NewError(err)
	}

	return resp, nil
}
