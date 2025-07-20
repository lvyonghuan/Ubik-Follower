package main

import (
	"Ubik-Follower/api"
	"Ubik-Follower/engine"
	"os"
	"os/exec"
)

const inTest = false

func main() {
	// Set the terminal to UTF-8 encoding
	cmd := exec.Command("chcp", "65001")
	cmd.Stdout = os.Stdout
	cmd.Run()

	uFollower := engine.InitEngine(inTest)

	uFollower.Log.Info("Ubik-Follower is starting...")
	initAPI(uFollower, inTest)
}

func initAPI(e *engine.UFollower, inTest bool) {
	err := api.InitAPI(e, inTest)
	if err != nil {
		e.Log.Fatal(err)
	}
}
