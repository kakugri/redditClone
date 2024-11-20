// internal/simulator/simulator.go
package simulator

import (
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/kakugri/redditClone/internal/engine"
)

type Simulator struct {
	enginePID *actor.PID
	users     []*SimulatedUser
	numUsers  int
	metrics   *Metrics
}

type SimulatedUser struct {
	userID    string
	connected bool
	pid       *actor.PID
}

type Metrics struct {
	TotalPosts    int64
	TotalComments int64
	TotalVotes    int64
	ActiveUsers   int64
	StartTime     time.Time
}

type UserActor struct {
	enginePID     *actor.PID
	postFrequency time.Duration
}

type NewUserActor struct {
	enginePID     *actor.PID
	postFrequency time.Duration
}

func (u *NewUserActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		// Simulate user behavior
		go u.simulateActivity(context)
	}
}

func NewSimulator(enginePID *actor.PID, numUsers int) *Simulator {
	return &Simulator{
		enginePID: enginePID,
		numUsers:  numUsers,
		metrics: &Metrics{
			StartTime: time.Now(),
		},
	}
}

// Zipf distribution helper
func generateZipfDistribution(s float64, v float64, n uint64) []uint64 {
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), s, v, n)
	distribution := make([]uint64, n)

	for i := uint64(0); i < n; i++ {
		distribution[i] = zipf.Uint64()
	}

	return distribution
}

func (s *Simulator) Start() {
	// Generate Zipf distribution for user activity
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.5, 1.0, uint64(s.numUsers))

	// Create user actors
	for i := 0; i < s.numUsers; i++ {
		activity := zipf.Uint64()
		props := actor.PropsFromProducer(func() actor.Actor {
			return &NewUserActor{
				s.enginePID,
				time.Duration(activity) * time.Second,
			}
		})

		pid := actor.NewActorSystem().Root.Spawn(props)
		s.users = append(s.users, &SimulatedUser{
			userID:    generateID(),
			connected: true,
			pid:       pid,
		})
	}
}

func (u *NewUserActor) simulateActivity(context actor.Context) {
	for {
		// Simulate connection/disconnection
		if rand.Float64() < 0.1 { // 10% chance to disconnect
			time.Sleep(time.Duration(rand.Intn(300)) * time.Second) // Sleep 0-5 minutes
			continue
		}

		// Simulate various actions
		switch rand.Intn(5) {
		case 0:
			// Create post
			context.Send(u.enginePID, &engine.CreatePostMsg{
				Title:   "Simulated Post",
				Content: "Content " + generateID(),
			})
		case 1:
			// Create comment
			context.Send(u.enginePID, &engine.CreateCommentMsg{
				Content: "Comment " + generateID(),
			})
		case 2:
			// Vote
			context.Send(u.enginePID, &engine.VoteMsg{
				IsUpvote: rand.Float64() < 0.7, // 70% chance of upvote
			})
		}

		time.Sleep(u.postFrequency)
	}
}

func (s *Simulator) GetMetrics() *Metrics {
	return s.metrics
}

func generateID() string {
	return time.Now().String()
}
