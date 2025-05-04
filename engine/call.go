package engine

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"

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

	//Marshal the node data
	jsonData, err := ujson.Marshal(nodes)
	if err != nil {
		engine.Log.Error(uerr.NewError(err))
		return
	}

	u := plugin.PluginMetaData.Uri + initNodes
	req, err := http.NewRequest("PUT", u, bytes.NewBuffer(jsonData))
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
			baseUrl, err := url.Parse(node.pluginInfo.PluginMetaData.Uri + setParams)
			if err != nil {
				engine.Log.Error(uerr.NewError(err))
				continue
			}

			//Set the node id to the url
			q := baseUrl.Query()
			q.Set("id", strconv.Itoa(node.ID))
			baseUrl.RawQuery = q.Encode()

			jsonData, err := ujson.Marshal(node.params)
			if err != nil {
				engine.Log.Error(uerr.NewError(err))
				continue
			}

			req, err := http.NewRequest("PUT", baseUrl.String(), bytes.NewBuffer(jsonData))
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
	u := plugin.PluginMetaData.Uri + run
	req, err := http.NewRequest("GET", u, nil)
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
