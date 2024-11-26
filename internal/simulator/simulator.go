// internal/simulator/simulator.go
package simulator

import (
	"log"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/kakugri/redditClone/internal/proto"
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

type NewUserActor struct {
	enginePID     *actor.PID
	postFrequency time.Duration
	userID        string
	subredditID   string
}

func (u *NewUserActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
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

func (s *Simulator) Start() {
	// Create user actors
	for i := 0; i < s.numUsers; i++ {
		userID := generateID()
		subredditID := "subreddit-1" // Simulate a fixed subreddit for now
		props := actor.PropsFromProducer(func() actor.Actor {
			return &NewUserActor{
				enginePID:     s.enginePID,
				postFrequency: time.Duration(rand.Intn(10)+1) * time.Second,
				userID:        userID,
				subredditID:   subredditID,
			}
		})

		pid := actor.NewActorSystem().Root.Spawn(props)
		s.users = append(s.users, &SimulatedUser{
			userID:    userID,
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
		switch rand.Intn(3) {
		case 0:
			// Simulate creating a post
			msg := &proto.CreatePostMsg{
				Title:       "Simulated Post Title",
				Content:     "Simulated Post Content",
				AuthorId:    u.userID,
				SubredditId: u.subredditID,
			}
			context.Send(u.enginePID, msg)
			log.Printf("User %s sent a post to subreddit %s", u.userID, u.subredditID)
		case 1:
			// Create comment
			msg := &proto.RegisterUserMsg{
				Username: u.userID,
			}
			context.Send(u.enginePID, msg)
			log.Printf("User %s created", u.userID)
		case 2:
			// Create comment
			msg := &proto.CreateSubredditMsg{
				Name:        u.subredditID,
				Description: "Simulated Subreddit",
				CreatorId:   u.userID,
			}
			// log.Printf("Sending CreateCommentMsg: %+v", msg)
			context.Send(u.enginePID, msg)
			log.Printf("User %s created subreddit %s", u.userID, u.subredditID)
		}

		time.Sleep(u.postFrequency)
	}
}

func (s *Simulator) GetMetrics() *Metrics {
	return s.metrics
}

func generateID() string {
	return time.Now().Format("20060102150405.000")
}
