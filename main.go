package main

import (
	"Ubik-Follower/api"
	"Ubik-Follower/engine"
)

func main() {
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
