package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get the list of plugins
func getPluginList(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	plugins := engine.GetPluginList()

	successResponse(c, plugins)
}

func addRuntimeNode(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	pluginName := c.Query("plugin_name")
	nodeName := c.Query("node_name")
	idString := c.Query("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
	}

	err = engine.NewRuntimeNode(pluginName, nodeName, id)
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "New runtime node created")
}

func deleteRuntimeNode(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	idString := c.Query("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
	}

	err = engine.DeleteRuntimeNode(id)
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "Runtime node deleted")
}

func updateEdge(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	producerIDString := c.Query("producer_id")
	consumerIDString := c.Query("consumer_id")
	producerPortName := c.Query("producer_port_name")
	consumerPortName := c.Query("consumer_port_name")
	uri := c.Query("uri")

	producerID, err := strconv.Atoi(producerIDString)
	if err != nil {
		errorResponse(c, 400, "producer_id err: "+err.Error())
	}

	consumerID, err := strconv.Atoi(consumerIDString)
	if err != nil {
		errorResponse(c, 400, "consumer_id err: "+err.Error())
	}

	err = engine.UpdateEdge(producerID, consumerID, producerPortName, consumerPortName, uri)
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "Edge updated")
}

func deleteEdge(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	producerIDString := c.Query("producer_id")
	consumerIDString := c.Query("consumer_id")
	producerPortName := c.Query("producer_port_name")
	consumerPortName := c.Query("consumer_port_name")

	producerID, err := strconv.Atoi(producerIDString)
	if err != nil {
		errorResponse(c, 400, "producer_id should be an integer")
	}

	consumerID, err := strconv.Atoi(consumerIDString)
	if err != nil {
		errorResponse(c, 400, "consumer_id should be an integer")
	}

	err = engine.DeleteEdge(producerID, consumerID, producerPortName, consumerPortName)
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "Edge deleted")
}

func putParams(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	idString := c.Query("id")
	paramsJson := c.PostForm("params")

	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
	}

	err = engine.PutParams(id, paramsJson)
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "Params updated")
}

func waitPrepare(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	// Wait for the prepare signal
	engine.WaitPrepare()

	successResponse(c, "Prepare checked")
}

func runPlugins(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	err = engine.RunPlugins()
	if err != nil {
		errorResponse(c, 400, err.Error())
	}

	successResponse(c, "Plugins running")
}
