package engine

import "testing"

func TestAddNode(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("example", "StartNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "ProcessNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "OutputNode", 3)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteNode(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("example", "StartNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "ProcessNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "OutputNode", 3)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.DeleteRuntimeNode(2)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateEdge(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("example", "StartNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "ProcessNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "OutputNode", 3)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(1, 2, "text", "inputText")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(1, 3, "text", "data")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(2, 3, "processedData", "data")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteEdge(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("example", "StartNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "ProcessNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("example", "OutputNode", 3)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(1, 2, "text", "inputText")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(1, 3, "text", "data")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(2, 3, "processedData", "data")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.DeleteEdge(1, 2, "text", "inputText")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.DeleteEdge(2, 3, "processedData", "data")
	if err != nil {
		t.Error(err)
	}
}
