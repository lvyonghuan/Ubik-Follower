package engine

import (
	"time"
)

//communicate with the leader

func (engine *UFollower) detectLeader() {
	for {
		if engine.Config.LeaderUrl != "" {

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

//func (engine *UFollower) findLeaderByURL() error {
//	url := engine.Config.LeaderUrl + "/follower" + "/init"
//}

func (engine *UFollower) broadCastToFindLeader() error {
	return nil
}

//func (engine *UFollower) postPlugins() error {
//	//get all plugin metadata
//	var plugins []uplugin.Plugin
//	for _, plugin := range engine.plugin.plugins {
//		plugins = append(plugins, *plugin.PluginMetaData)
//	}
//
//	//marshal the plugin metadata into JSON
//	jsonData, err := ujson.Marshal(plugins)
//	if err != nil {
//		return err
//	}
//}
