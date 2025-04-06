package main

import (
	"Ubik-Follower/api"
	"Ubik-Follower/engine"
	"os"
	"os/exec"
)

func main() {
	// Set the terminal to UTF-8 encoding
	cmd := exec.Command("chcp", "65001")
	cmd.Stdout = os.Stdout
	cmd.Run()

	uFollower := engine.InitEngine(false)

	uFollower.Log.Info("Ubik-Follower is starting...")
	initAPI(uFollower)
}

func initAPI(e *engine.UFollower) {
	err := api.InitAPI(e)
	if err != nil {
		e.Log.Fatal(err)
	}
}
