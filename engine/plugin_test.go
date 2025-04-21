package engine

import "testing"

func TestRunPlugin(t *testing.T) {
	e := InitEngine(true)

	err := e.NewRuntimeNode("AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("AddNum", "sumNode", 3)
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("AddNum", "selfIncreasingNode", 4)
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge(1, 3, "num_a", "num_a")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge(1, 2, "num_a", "input")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge(2, 3, "num_b", "num_b")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge(1, 4, "cycle_num", "input")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge(4, 1, "num_b", "current_cycle_num")
	if err != nil {
		t.Error(err)
	}
}
