// cmd/engine/main.go
package main

import (
    "log"
    "github.com/asynkron/protoactor-go/actor"
    "github.com/kakugri/redditClone/internal/engine"
)

func main() {
    system := actor.NewActorSystem()
    
    props := actor.PropsFromProducer(func() actor.Actor {
        return engine.NewRedditEngine()
    })

    pid := system.Root.Spawn(props)
    log.Printf("Reddit Engine started with PID: %v", pid)
    
    // Keep the engine running
    select {}
}