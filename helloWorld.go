package main

// Import the actor package
import (
	"fmt"
	"time"
	"github.com/asynkron/protoactor-go/actor"
)

// Define the Hello message with a Who property
type Hello struct {
	Who string
}

// Define the HelloWorldActor struct
type HelloWorldActor struct{}

// Implement the Receive method for message processing
func (state *HelloWorldActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Hello:
		// Print the greeting to the console
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	// Create a new ActorSystem
	system := actor.NewActorSystem()

	// Create a Props object for the HelloWorldActor
	props := actor.PropsFromProducer(func() actor.Actor { return &HelloWorldActor{} })

	// Spawn a new instance of the HelloWorldActor
	pid := system.Root.Spawn(props)

	system.Root.Send(pid, &Hello{Who: "World"})
	time.Sleep(1 * time.Second)
}
