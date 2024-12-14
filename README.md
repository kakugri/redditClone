# Reddit Clone Engine and Simulator
## Team Members:
* Kingdom Mutala Akugri
* Siddhant Kalgutkar
## Overview
This project implements a Reddit-like engine and a simulator to test its performance and scalability. The engine supports user registration, subreddit management, posting, commenting, and voting. The simulator mimics user interactions to generate meaningful metrics.
### Instructions to Run:
1) Clone the repo and open:
```sh
git clone git@github.com:kakugri/redditClone.git
cd redditClone
```
2) Run Engine in terminal using:
```sh
go run cmd/engine/redditEngine.go
```
Skip to 4 if running testing the REST API

3) Run Simulator in separate terminal connects to the engine and generates activity. Metrics are logged every minute.:
```sh
go run cmd/simulator/redditSimulator.go
```
Only run the next few steps if testing the REST API

4) Run API in separate terminal connects to the engine and generates activity. Metrics are logged every minute.:
```sh
go run cmd/api2/api2.go 
```
5) Run client simulator in separate terminal connects to the engine and generates activity. Metrics are logged every minute.:
```sh
go run cmd/clients/clients.go 
```
Features:
## Engine
* Register and manage user accounts.
* Create, join, and leave subreddits.
* Post, comment, and upvote/downvote content.
* Send and receive direct messages.
* Compute and track metrics like karma.
## Simulator
* Mimics thousands of user interactions.
* Models disconnection/reconnection behavior.
* Generates performance metrics.
Directory Structure:
* internal/engine: Core engine logic and models.
* internal/simulator: Simulator logic for user actions.
* cmd/engine: Entry point for the engine.
* cmd/simulator: Entry point for the simulator.
* internal/proto: Protobuf definitions for communication.
Metrics
#### The engine periodically reports metrics, including:
* Total posts
* Total comments
* Total votes
* Active users
* Total messages
Logs are available in the terminal during runtime.
Largest Network Size Tested:
* The implementation was successfully tested with a maximum of 100,000 nodes.