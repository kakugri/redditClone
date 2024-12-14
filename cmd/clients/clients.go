package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

const registerURL = "http://localhost:8081/api/register"
const postURL = "http://localhost:8081/api/posts"
const subredditURL = "http://localhost:8081/api/newsubreddit"
const commentURL = "http://localhost:8081/api/comment"

func main() {
	var wg sync.WaitGroup
	clientCount := 100000

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			username := "user" + strconv.Itoa(clientID)
			simulateRequest(username)
		}(i)
	}

	wg.Wait()
	log.Println("All clients finished")
}

func simulateRequest(username string) {
	switch rand.Intn(4) {
	case 0:
		data := map[string]string{"username": username}
		body, _ := json.Marshal(data)

		resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error registering user %s: %v", username, err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Registered user: %s, Status: %d", username, resp.StatusCode)
	case 1:
		title := "Simulated Post Title"
		content := "Simulated Post Content"
		data := map[string]string{"title": title, "content": content}
		body, _ := json.Marshal(data)

		resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error sending post with title %s: %v", title, err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Sent post by user: %s, Status: %d", username, resp.StatusCode)
	case 2:
		description := "Simulated Subreddit Description"
		creatorid := username
		data := map[string]string{"username": username, "description": description, "creatorid": creatorid}
		body, _ := json.Marshal(data)

		resp, err := http.Post(subredditURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error creating subreddit with description %s: %v", description, err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Created subreddit by user: %s, Status: %d", username, resp.StatusCode)
	case 3:
		content := "Simulated Comment Content"
		data := map[string]string{"content": content}
		body, _ := json.Marshal(data)

		resp, err := http.Post(commentURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error creating comment by user %s: %v", username, err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Created comment by user: %s, Status: %d", username, resp.StatusCode)
	}
}
