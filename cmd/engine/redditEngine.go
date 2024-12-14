// cmd/engine/redditEngine.go
package main

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/kakugri/redditClone/internal/engine"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("localhost", 8080)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	props := actor.PropsFromProducer(func() actor.Actor {
		return engine.NewRedditEngine()
	})

	pid, err := system.Root.SpawnNamed(props, "reddit-engine")
	if err != nil {
		log.Fatalf("Failed to start engine: %v", err)
	}
	log.Printf("RedditEngine started: Address=%s, Id=%s", pid.GetAddress(), pid.GetId())

	select {}
}
