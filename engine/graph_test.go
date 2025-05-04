package engine

import "testing"

func TestAddNode(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "sumNode", 3)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteNode(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "sumNode", 3)
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

	err := uFollower.NewRuntimeNode("AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "sumNode", 3)
	if err != nil {
		t.Error(err)
	}

	// Connect the num_a output of startNode to the num_a input of sumNode
	err = uFollower.UpdateEdge(1, 3, "num_a", "num_a", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	// Connect the num_a output of startNode to the input of selfIncreasingNode
	err = uFollower.UpdateEdge(1, 2, "num_a", "input", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	// Connect the num_b output of selfIncreasingNode to the num_b input of sumNode
	err = uFollower.UpdateEdge(2, 3, "num_b", "num_b", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	// Just test: connect the num_a output of startNode to the num_b input of sumNode
	err = uFollower.UpdateEdge(1, 3, "num_a", "num_b", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteEdge(t *testing.T) {
	uFollower := InitEngine(true)

	err := uFollower.NewRuntimeNode("AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = uFollower.NewRuntimeNode("AddNum", "sumNode", 3)
	if err != nil {
		t.Error(err)
	}

	// 添加边
	err = uFollower.UpdateEdge(1, 3, "num_a", "num_a", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(1, 2, "num_a", "input", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.UpdateEdge(2, 3, "num_b", "num_b", "http://localhost:8080/input")
	if err != nil {
		t.Error(err)
	}

	// 删除边
	err = uFollower.DeleteEdge(1, 2, "num_a", "input")
	if err != nil {
		t.Error(err)
	}

	err = uFollower.DeleteEdge(2, 3, "num_b", "num_b")
	if err != nil {
		t.Error(err)
	}
}
