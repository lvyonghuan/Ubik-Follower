package engine_test

import (
	"Ubik-Follower/api"
	"Ubik-Follower/engine"
	"testing"
	"time"

	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/uplugin"
)

func TestRunPlugin(t *testing.T) {
	e := engine.InitEngine(true)
	go api.InitAPI(e, true)
	time.Sleep(5 * time.Second)

	err := e.CreateWorkflow("testWorkflow")
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("testWorkflow", "AddNum", "startNode", 1)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)

	err = e.NewRuntimeNode("testWorkflow", "AddNum", "selfIncreasingNode", 2)
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("testWorkflow", "AddNum", "sumNode", 3)
	if err != nil {
		t.Error(err)
	}

	err = e.NewRuntimeNode("testWorkflow", "AddNum", "selfIncreasingNode", 4)
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 1, 3, "num_a", "num_a", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 1, 2, "num_a", "input", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 2, 3, "num_b", "num_b", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 1, 4, "cycle_num", "input", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 4, 1, "num_b", "current_cycle_num", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	err = e.UpdateEdge("testWorkflow", 3, 1, "sum", "num_a", "http://localhost:14535")
	if err != nil {
		t.Error(err)
	}

	params := make(uplugin.Params)
	params["init_num"] = 0
	params["cycle_num"] = 10
	paramsJson, err := ujson.Marshal(params)
	if err != nil {
		t.Error(err)
	}

	err = e.PutParams("testWorkflow", 1, paramsJson)
	if err != nil {
		t.Error(err)
	}

	e.InitPluginsNodes()

	e.WaitPrepare()
	err = e.RunPlugins()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(10 * time.Second)
}
