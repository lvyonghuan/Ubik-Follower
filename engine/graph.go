package engine

import (
	"errors"
	"strconv"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

//FIXME:必须考虑consumer node不在该follower辖域内的问题。对consumer的check需要放弃。

//Directed graph of the workflow describing the connection status of nodes
//under the current follower's control.
//The graph consists of a series of runtime nodes and directed edges.

// NewRuntimeNode creates a new runtime node.
func (engine *UFollower) NewRuntimeNode(workflowName, pluginName, nodeName string, id int) error {
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
		OutputEdges: make(map[string][]Edge),
	}

	//Add the runtime node to the follower's runtime node map
	//Check if the id already exists and if the workflow exists
	if workflow, isExist := engine.workflows[workflowName]; !isExist {
		return uerr.NewError(errors.New("workflow " + workflowName + " not exist"))
	} else {
		if _, isExist := workflow[id]; isExist {
			return uerr.NewError(errors.New("runtime node with id " + strconv.Itoa(id) + " already exists in workflow " + workflowName))
		}
	}
	engine.workflows[workflowName][id] = runtimeNode
	engine.plugin.plugins[pluginName].runtimeNodes[id] = runtimeNode //Mount node under the plugin

	//Check if the plugin is already mounted
	//If not, then will mount the plugin
	go engine.checkMount(plugin)

	engine.Log.Debug("New runtime node " + nodeName + " created, id: " + strconv.Itoa(id))

	return nil
}

// DeleteRuntimeNode deletes a runtime node.
func (engine *UFollower) DeleteRuntimeNode(workflowName string, id int) error {
	if workflow, isExist := engine.workflows[workflowName]; !isExist {
		return uerr.NewError(errors.New("workflow " + workflowName + " not exist"))
	} else {
		if runtimeNode, isExist := workflow[id]; !isExist {
			return uerr.NewError(errors.New("runtime node with id " + strconv.Itoa(id) + " not exist in workflow " + workflowName))
		} else {
			//delete the runtime node from the workflow
			for portName, edges := range runtimeNode.OutputEdges {
				for _, edge := range edges {
					err := engine.deleteEdge(workflow, id, edge.ConsumerID, portName, edge.ConsumerPortName)
					if err != nil {
						engine.Log.Error(err)
						continue
					}
				}
			}

			//delete the runtime node from the plugin's runtime node map
			delete(runtimeNode.pluginInfo.runtimeNodes, id)
		}
	}

	//delete the runtime node from the follower's runtime node map
	delete(engine.workflows[workflowName], id)

	//TODO design the unmount plugin logic. Not just ==0 unmount
	//go engine.checkUnmount(engine.runtimeNodes[id].pluginInfo)

	engine.Log.Debug("Runtime node " + strconv.Itoa(id) + " deleted")

	return nil
}

// UpdateEdge updates the edge between two runtime nodes.
func (engine *UFollower) UpdateEdge(workflowName string, producerID, consumerID int, producerPortName, consumerPortName, uri string) error {
	workflow, isExist := engine.workflows[workflowName]
	if !isExist {
		return uerr.NewError(errors.New("workflow " + workflowName + " not exist"))
	}

	producer, isExist := workflow[producerID]
	if !isExist {
		return uerr.NewError(errors.New("producer runtime node " + strconv.Itoa(producerID) + " not exist"))
	}

	edge := Edge{
		ProducerID:       producerID,
		ConsumerID:       consumerID,
		ProducerPortName: producerPortName,
		ConsumerPortName: consumerPortName,
		channel:          make(chan []byte),
		Addr:             uri,
	}

	producer.OutputEdges[producerPortName] = append(producer.OutputEdges[producerPortName], edge)

	engine.Log.Debug("Edge updated between runtime node " + strconv.Itoa(producerID) + " and " + strconv.Itoa(consumerID))

	return nil
}

// DeleteEdge deletes the edge between two runtime nodes.
func (engine *UFollower) DeleteEdge(workflowName string, producerID, consumerID int, producerPortName, consumerPortName string) error {
	workflow, isExist := engine.workflows[workflowName]
	if !isExist {
		return uerr.NewError(errors.New("workflow " + workflowName + " not exist"))
	}

	err := engine.deleteEdge(workflow, producerID, consumerID, producerPortName, consumerPortName)
	if err != nil {
		return err
	}

	return nil
}

func (engine *UFollower) deleteEdge(workflow runtimeNodes, producerID, consumerID int, producerPortName, consumerPortName string) error {
	producer, isExist := workflow[producerID]
	if !isExist {
		return uerr.NewError(errors.New("producer runtime node " + strconv.Itoa(producerID) + " not exist"))
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

	engine.Log.Debug("Edge deleted between runtime node " + strconv.Itoa(producerID) + " and " + strconv.Itoa(consumerID))
	return nil
}

// DescriptionGraph describes the graph of the workflow.
func DescriptionGraph() {

}
