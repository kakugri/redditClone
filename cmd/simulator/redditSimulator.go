// cmd/simulator/main.go
package main

import (
	"log"
	"runtime"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/kakugri/redditClone/internal/proto"
	"github.com/kakugri/redditClone/internal/simulator"
	"github.com/shirou/gopsutil/v3/cpu"
)

func main() {
	// Initialize the actor system
	system := actor.NewActorSystem()

	// Configure the remote actor system
	config := remote.Configure("localhost", 0) // The simulator listens on a dynamic port
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Target the RedditEngine actor
	enginePID := actor.NewPID("127.0.0.1:8080", "reddit-engine")
	log.Printf("Simulator targeting engine at Address=%s, Id=%s", enginePID.GetAddress(), enginePID.GetId())

	// Subscribe to DeadLetter events for debugging
	system.EventStream.Subscribe(func(evt interface{}) {
		if deadLetter, ok := evt.(actor.DeadLetterEvent); ok {
			log.Printf("DeadLetter Event - PID: %+v, Message: %+v, Sender: %+v",
				deadLetter.PID, deadLetter.Message, deadLetter.Sender)
		}
	})

	// Delay startup to ensure the engine is initialized
	log.Println("Waiting for engine to initialize...")
	time.Sleep(10 * time.Second)

	// Send a test CreatePostMsg to the engine
	log.Println("Sending test CreatePostMsg to engine...")
	system.Root.Send(enginePID, &proto.CreatePostMsg{
		Title:       "Test Post",
		Content:     "Test Content",
		AuthorId:    "test-author-id",
		SubredditId: "test-subreddit-id",
	})
	log.Println("Test CreatePostMsg sent.")
	time.Sleep(20 * time.Second)

	// Start the simulator with 5 simulated users
	sim := simulator.NewSimulator(enginePID, 2100)
	sim.Start()
	log.Println("Simulator started.")

	// Print metrics every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			metrics := sim.GetMetrics()
			// Capture memory usage
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Capture CPU usage
			cpuPercent, _ := cpu.Percent(0, false)
			log.Printf("Metrics Report: TotalPosts=%d, ActiveUsers=%d, TotalVotes=%d, TotalComments=%d, TotalMessages=%d, Memory=%.2f MB, CPU=%.2f%%", metrics.TotalPosts, metrics.ActiveUsers,
				metrics.TotalVotes, metrics.TotalComments, metrics.TotalMessages, float64(m.Alloc)/1024/1024,
				cpuPercent[0])
		}
	}()

	// Keep the simulator process alive
	select {}
}
