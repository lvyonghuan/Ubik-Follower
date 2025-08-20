package engine

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

//communicate with the leader

func (engine *UFollower) detectLeader() {
	for {
		if engine.Config.LeaderAddr != "" {
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
	url := engine.Config.LeaderAddr + "/follower" + "/init"

	// Prepare the request URL with UUID
	reqURL := url + "?UUID=" + engine.UUID + "&Addr=" + engine.Config.IP + ":" + engine.Config.Port

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
	url := engine.Config.LeaderAddr + "/follower" + "/list"
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

// Heartbeat is a struct that manages the heartbeat mechanism for the follower
type heartbeat struct {
	conn     *net.UDPConn
	interval time.Duration
}

type heartbeatPacket struct {
	UUID string `json:"UUID"`
}

func (engine *UFollower) initHeartbeat() error {
	engine.heartbeat = &heartbeat{}

	// Initialize the heartbeat sender
	engine.heartbeat.interval = time.Duration(engine.Config.HeartBeatInterval) * time.Second

	// Start sending heartbeats to the leader
	err := engine.startHeartbeat()
	if err != nil {
		return uerr.NewError(err)
	}

	engine.Log.Debug("Heartbeat sender initialized")
	return nil
}

func (engine *UFollower) startHeartbeat() error {
	parsedURL, err := url.Parse(engine.Config.LeaderAddr)
	if err != nil {
		return uerr.NewError(err)
	}

	// Create an UDP connection
	addr, err := net.ResolveUDPAddr("udp", ":"+parsedURL.Port())
	if err != nil {
		return uerr.NewError(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return uerr.NewError(err)
	}
	engine.heartbeat.conn = conn

	// Start sending heartbeats at regular intervals
	ticker := time.NewTicker(engine.heartbeat.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := engine.sendHeartbeat()
				if err != nil {
					engine.Log.Error(uerr.NewError(err))
				}
			}
		}
	}()

	return nil
}

func (engine *UFollower) sendHeartbeat() error {
	// Create a heartbeat packet
	packet := heartbeatPacket{
		UUID: engine.UUID,
	}

	// Marshal the packet to JSON
	jsonData, err := ujson.Marshal(packet)
	if err != nil {
		return uerr.NewError(err)
	}

	// Send the heartbeat packet to the leader
	_, err = engine.heartbeat.conn.Write(jsonData)
	if err != nil {
		return uerr.NewError(err)
	}

	engine.Log.Debug("Heartbeat sent to leader")
	return nil
}
