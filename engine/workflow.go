package engine

import (
	"errors"

	"github.com/lvyonghuan/Ubik-Util/uerr"
)

func (engine *UFollower) CreateWorkflow(workflowName string) error {
	// If the workflow already exists, return an error
	if _, exists := engine.workflows[workflowName]; exists {
		return uerr.NewError(errors.New("workflow " + workflowName + " already exists"))
	}

	// Create a new workflow and add it to the workflows map
	engine.workflows[workflowName] = make(runtimeNodes)
	return nil
}

func (engine *UFollower) DeleteWorkflow(workflowName string) error {
	// If the workflow does not exist, return an error
	if _, exists := engine.workflows[workflowName]; !exists {
		return uerr.NewError(errors.New("workflow " + workflowName + " does not exist"))
	}

	// Delete the workflow from the workflows map
	delete(engine.workflows, workflowName)
	return nil
}
