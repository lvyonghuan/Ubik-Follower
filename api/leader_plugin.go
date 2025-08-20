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

func createWorkflow(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	workflowName := c.Query("workflow_name")
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.CreateWorkflow(workflowName)
	if err != nil {
		errorResponse(c, 400, "create workflow err: "+err.Error())
		return
	}

	successResponse(c, "Workflow created successfully")
}

// DeleteWorkflow deletes a workflow by its name
func deleteWorkflow(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
		return
	}

	workflowName := c.Query("workflow_name")
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.DeleteWorkflow(workflowName)
	if err != nil {
		errorResponse(c, 400, "delete workflow err: "+err.Error())
		return
	}

	successResponse(c, "Workflow deleted successfully")
}

func addRuntimeNode(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	pluginName := c.Query("plugin_name")
	nodeName := c.Query("node_name")
	idString := c.Query("id")
	workflowName := c.Query("workflow_name")

	// Check validation
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
		return
	}
	if pluginName == "" {
		errorResponse(c, 400, "plugin_name is required, but get empty")
		return
	}
	if nodeName == "" {
		errorResponse(c, 400, "node_name is required, but get empty")
		return
	}
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.NewRuntimeNode(workflowName, pluginName, nodeName, id)
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}

	successResponse(c, "New runtime node created")
}

func deleteRuntimeNode(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	idString := c.Query("id")
	workflowName := c.Query("workflow_name")

	// Check validation
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
		return
	}
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.DeleteRuntimeNode(workflowName, id)
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
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
	addr := c.Query("addr")
	workflowName := c.Query("workflow_name")

	// Check validation
	producerID, err := strconv.Atoi(producerIDString)
	if err != nil {
		errorResponse(c, 400, "producer_id err: "+err.Error())
		return
	}
	consumerID, err := strconv.Atoi(consumerIDString)
	if err != nil {
		errorResponse(c, 400, "consumer_id err: "+err.Error())
		return
	}
	if producerPortName == "" {
		errorResponse(c, 400, "producer_port_name is required, but get empty")
		return
	}
	if consumerPortName == "" {
		errorResponse(c, 400, "consumer_port_name is required, but get empty")
		return
	}
	if addr == "" {
		errorResponse(c, 400, "addr is required, but get empty")
		return
	}
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.UpdateEdge(workflowName, producerID, consumerID, producerPortName, consumerPortName, addr)
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
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
	workflowName := c.Query("workflow_name")

	// Check validation
	producerID, err := strconv.Atoi(producerIDString)
	if err != nil {
		errorResponse(c, 400, "producer_id should be an integer")
		return
	}
	consumerID, err := strconv.Atoi(consumerIDString)
	if err != nil {
		errorResponse(c, 400, "consumer_id should be an integer")
		return
	}
	if producerPortName == "" {
		errorResponse(c, 400, "producer_port_name is required, but get empty")
		return
	}
	if consumerPortName == "" {
		errorResponse(c, 400, "consumer_port_name is required, but get empty")
		return
	}
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.DeleteEdge(workflowName, producerID, consumerID, producerPortName, consumerPortName)
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}

	successResponse(c, "Edge deleted")
}

func putParams(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}

	idString := c.Query("id")
	workflowName := c.Query("workflow_name")

	body, err := c.GetRawData()
	if err != nil {
		errorResponse(c, 400, "read body err: "+err.Error())
		return
	}

	// Check validation
	id, err := strconv.Atoi(idString)
	if err != nil {
		errorResponse(c, 400, "id err: "+err.Error())
		return
	}
	if workflowName == "" {
		errorResponse(c, 400, "workflow_name is required, but get empty")
		return
	}

	err = engine.PutParams(workflowName, id, body)
	if err != nil {
		errorResponse(c, 400, err.Error())
		return
	}

	successResponse(c, "Params updated")
}

func waitPrepare(c *gin.Context) {
	engine, err := getEngine(c)
	if err != nil {
		fatalErrHandel(c, err)
	}
	// Init the Plugin nodes
	engine.InitPluginsNodes()

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
		return
	}

	successResponse(c, "Plugins running")
}
