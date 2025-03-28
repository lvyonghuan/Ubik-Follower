package engine

import "testing"

func TestInit(t *testing.T) {
	engine := InitEngine(true)
	if engine == nil {
		t.Error("engine is nil")
	}
}
