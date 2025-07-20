package api

import (
	"Ubik-Follower/engine"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lvyonghuan/Ubik-Util/uerr"
)

const (
	engineKey = "engine" //key for the engine in the gin context
)

//API for the leader
//TODO: Authenticate the leader by middleware

// InitAPI initializes the API
func InitAPI(engine *engine.UFollower, inTest bool) error {
	r := gin.Default()
	// Set the mode based on whether it's in test or not
	if !inTest {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine.SetGinEngine(r)

	//set the engine in the context
	r.Use(func(c *gin.Context) {
		c.Set(engineKey, engine)
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/list", getPluginList) //get a list of follower's installed plugins
		api.GET("/info")                //get metadata of a plugin
		api.PUT("/download")            //Download a plugin. Three steps: 1. download 2. unzip 3.init

		node := api.Group("/node")
		{
			node.PUT("/", addRuntimeNode)       //create a runtime node
			node.DELETE("/", deleteRuntimeNode) //delete a runtime node
			node.PUT("/edge", updateEdge)       //update a edge from nodeA(producer) to nodeB(consumer)
			node.DELETE("/edge", deleteEdge)    //delete a edge from nodeA(producer) to nodeB(consumer)
			node.PUT("/param", putParams)       //update the parameter of a node
		}

		api.GET("/prepare", waitPrepare) //a signal to tell follower prepare to run, response when all plugins are mounted and no error
		api.GET("/run", runPlugins)      //run the follower
		api.GET("/stop")                 //stop the running follower
	}

	plugin := r.Group("/plugin")
	{
		plugin.GET("/init", pluginInit) //when plugin in running, plugin should call this api
		plugin.GET("/status")           //report the status of the plugin
		plugin.POST("/error")           //report the error of the plugin
		plugin.PUT("/info")             //report the info of the plugin
		plugin.GET("/ping")             //heartbeat

		//TODO: use follower to forward the output of the node
		plugin.PUT("/output")
	}

	engine.Log.Info("Ubik-Follower running now, listening on " + engine.Config.Port)

	err := r.Run(engine.Config.IP + ":" + engine.Config.Port)
	if err != nil {
		return uerr.NewError(err)
	}
	return nil
}

// retrieves the engine from the context
func getEngine(c *gin.Context) (*engine.UFollower, error) {
	engineVal, isExist := c.Get(engineKey)
	if !isExist {
		return nil, uerr.NewError(errors.New("engine not exist in the context"))
	}
	return engineVal.(*engine.UFollower), nil
}
