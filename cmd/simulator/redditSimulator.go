// cmd/simulator/main.go
package main

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/kakugri/redditClone/internal/simulator"
)

func main() {
	// system := actor.NewActorSystem()
	enginePID := actor.NewPID("localhost:8080", "reddit-engine")

	sim := simulator.NewSimulator(enginePID, 1000) // Simulate 1000 users
	sim.Start()

	// Print metrics every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			metrics := sim.GetMetrics()
			log.Printf("Metrics: %+v", metrics)
		}
	}()

	// Keep the simulator running
	select {}
}
