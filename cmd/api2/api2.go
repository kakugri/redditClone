package main

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/kakugri/redditClone/internal/api2"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("localhost", 0)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Target the RedditEngine actor
	enginePID := actor.NewPID("127.0.0.1:8080", "reddit-engine")
	log.Printf("Simulator targeting engine at Address=%s, Id=%s", enginePID.GetAddress(), enginePID.GetId())

	router := api2.SetupRouter(system, enginePID)

	log.Println("REST API server running on :8081")
	router.Run(":8081")
}
