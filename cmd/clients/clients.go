package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const apiURL = "http://localhost:8081/api/register"

func main() {
	var wg sync.WaitGroup
	clientCount := 100

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			username := "user" + strconv.Itoa(clientID)
			registerUser(username)
		}(i)
	}

	wg.Wait()
	log.Println("All clients finished")
}

func registerUser(username string) {
	data := map[string]string{"username": username}
	body, _ := json.Marshal(data)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error registering user %s: %v", username, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Registered user: %s, Status: %d", username, resp.StatusCode)
}
