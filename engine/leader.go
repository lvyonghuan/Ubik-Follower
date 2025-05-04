package engine

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

//communicate with the leader

func (engine *UFollower) detectLeader() {
	for {
		if engine.Config.LeaderUrl != "" {
			err := engine.findLeaderByURL()
			if err != nil {
				engine.Log.Error(err)
				time.Sleep(5 * time.Second)
				continue //retry
			}
		} else { //broadcast to find leader
			err := engine.broadCastToFindLeader()
			if err != nil {
				engine.Log.Error(err)
				time.Sleep(5 * time.Second)
				continue //retry
			}
		}

		break
	}
}

func (engine *UFollower) findLeaderByURL() error {
	url := engine.Config.LeaderUrl + "/follower" + "/init"

	// Prepare the request URL with UUID
	reqURL := url + "?UUID=" + engine.UUID

	// Send a GET request to the leader
	resp, err := http.Get(reqURL)
	if err != nil {
		return uerr.NewError(fmt.Errorf("failed to connect to leader: %v", err))
	}
	defer resp.Body.Close()

	// Check the response status code
	//TODO:检查消息体的状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return uerr.NewError(fmt.Errorf("leader responded with status code %d: %s", resp.StatusCode, string(body)))
	}

	// Successfully connected to the leader
	engine.Log.Info("Successfully connected to leader")
	return nil
}

func (engine *UFollower) broadCastToFindLeader() error {
	return nil
}

func (engine *UFollower) postPlugins() error {
	url := engine.Config.LeaderUrl + "/follower" + "/list"
	reqURL := url + "?UUID=" + engine.UUID

	//get all plugin metadata
	plugins := make(map[string]*uplugin.Plugin)
	for _, plugin := range engine.plugin.plugins {
		plugins[plugin.PluginMetaData.Name] = plugin.PluginMetaData
	}

	//marshal the plugin metadata into JSON
	jsonData, err := ujson.Marshal(plugins)
	if err != nil {
		return err
	}

	//send the JSON to the leader
	resp, err := http.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return uerr.NewError(fmt.Errorf("failed to connect to leader: %v", err))
	}

	defer resp.Body.Close()
	// Check the response status code
	//TODO:检查消息体的状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return uerr.NewError(fmt.Errorf("leader responded with status code %d: %s", resp.StatusCode, string(body)))
	}

	// Successfully post the plugin list to the leader
	engine.Log.Info("Successfully post the plugin list to leader")
	return nil
}
