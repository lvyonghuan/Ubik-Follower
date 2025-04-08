package engine

import (
	"errors"
	"strconv"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

//Directed graph of the workflow describing the connection status of nodes
//under the current follower's control.
//The graph consists of a series of runtime nodes and directed edges.

// NewRuntimeNode creates a new runtime node.
func (engine *UFollower) NewRuntimeNode(pluginName, nodeName string, id int) error {
	plugin, isExist := engine.plugin.plugins[pluginName]
	if !isExist {
		return uerr.NewError(errors.New("plugin " + pluginName + " not exist"))
	}

	if _, isExist := plugin.PluginMetaData.Nodes[nodeName]; !isExist {
		return uerr.NewError(errors.New("node " + nodeName + " not exist"))
	}

	runtimeNode := &RuntimeNode{
		ID:          id,
		NodeName:    nodeName,
		pluginInfo:  plugin,
		params:      make(uplugin.Params),
		inputEdges:  make(map[string][]Edge),
		OutputEdges: make(map[string][]Edge),
	}

	//Add the runtime node to the follower's runtime node map
	//Check if the id already exists
	if _, isExist := engine.runtimeNodes[id]; isExist {
		return uerr.NewError(errors.New("runtime node id " + strconv.Itoa(id) + " already exist"))
	}
	engine.runtimeNodes[id] = runtimeNode
	engine.plugin.plugins[pluginName].runtimeNodes[id] = runtimeNode //Mount node under the plugin

	//Check if the plugin is already mounted
	//If not, then will mount the plugin
	go engine.checkMount(plugin)

	engine.Log.Debug("New runtime node " + nodeName + " created, id: " + strconv.Itoa(id))

	return nil
}

// DeleteRuntimeNode deletes a runtime node.
func (engine *UFollower) DeleteRuntimeNode(id int) error {
	if runtimeNode, isExist := engine.runtimeNodes[id]; !isExist {
		return uerr.NewError(errors.New("runtime node" + strconv.Itoa(id) + "not exist"))
	} else {
		//delete input edges
		for portName, edges := range runtimeNode.inputEdges {
			for _, edge := range edges {
				err := engine.DeleteEdge(edge.ProducerID, id, edge.ProducerPortName, portName)
				if err != nil {
					engine.Log.Error(err)
					continue
				}
			}
		}

		//delete output edges
		for portName, edges := range runtimeNode.OutputEdges {
			for _, edge := range edges {
				err := engine.DeleteEdge(id, edge.ConsumerID, portName, edge.ConsumerPortName)
				if err != nil {
					engine.Log.Error(err)
					continue
				}
			}
		}

		//delete the runtime node from the plugin's runtime node map
		delete(runtimeNode.pluginInfo.runtimeNodes, id)
	}

	//delete the runtime node from the follower's runtime node map
	delete(engine.runtimeNodes, id)

	//TODO design the unmount plugin logic. Not just ==0 unmount
	//go engine.checkUnmount(engine.runtimeNodes[id].pluginInfo)

	engine.Log.Debug("Runtime node " + strconv.Itoa(id) + " deleted")

	return nil
}

// UpdateEdge updates the edge between two runtime nodes.
func (engine *UFollower) UpdateEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	producer, isExist := engine.runtimeNodes[producerID]
	if !isExist {
		return uerr.NewError(errors.New("producer runtime node " + strconv.Itoa(producerID) + " not exist"))
	}

	consumer, isExist := engine.runtimeNodes[consumerID]
	if !isExist {
		return uerr.NewError(errors.New("consumer runtime node " + strconv.Itoa(consumerID) + " not exist"))
	}

	edge := Edge{
		ProducerID:       producerID,
		ConsumerID:       consumerID,
		ProducerPortName: producerPortName,
		ConsumerPortName: consumerPortName,
		channel:          make(chan []byte),
		Uri:              consumer.pluginInfo.PluginMetaData.Uri,
	}

	producer.OutputEdges[producerPortName] = append(producer.OutputEdges[producerPortName], edge)
	consumer.inputEdges[consumerPortName] = append(consumer.inputEdges[consumerPortName], edge)

	engine.Log.Debug("Edge updated between runtime node " + strconv.Itoa(producerID) + " and " + strconv.Itoa(consumerID))

	return nil
}

// DeleteEdge deletes the edge between two runtime nodes.
func (engine *UFollower) DeleteEdge(producerID, consumerID int, producerPortName, consumerPortName string) error {
	producer, isExist := engine.runtimeNodes[producerID]
	if !isExist {
		return uerr.NewError(errors.New("producer runtime node " + strconv.Itoa(producerID) + " not exist"))
	}

	consumer, isExist := engine.runtimeNodes[consumerID]
	if !isExist {
		return uerr.NewError(errors.New("consumer runtime node " + strconv.Itoa(consumerID) + " not exist"))
	}

	//Delete the edge from the producer's output edges
	//Two points can determine a straight line,
	//while ID and port name can determine a point
	edges := producer.OutputEdges[producerPortName]
	for i, edge := range edges {
		//Check if the edge is the one to be deleted
		if edge.ConsumerPortName == consumerPortName && edge.ConsumerID == consumerID {
			producer.OutputEdges[producerPortName] = append(edges[:i], edges[i+1:]...)
			break
		}
	}

	//Delete the edge from the consumer's input edges. Same as above.
	edges = consumer.inputEdges[consumerPortName]
	for i, edge := range edges {
		//Check if the edge is the one to be deleted
		if edge.ProducerPortName == producerPortName && edge.ProducerID == producerID {
			consumer.inputEdges[consumerPortName] = append(edges[:i], edges[i+1:]...)
			break
		}
	}

	engine.Log.Debug("Edge deleted between runtime node " + strconv.Itoa(producerID) + " and " + strconv.Itoa(consumerID))

	return nil
}

// DescriptionGraph describes the graph of the workflow.
func DescriptionGraph() {

}
