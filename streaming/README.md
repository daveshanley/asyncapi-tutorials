# AsyncAPI Tutorials for Golang: Streaming

This code accompanies an article on [quobix.com](https://quobix.com). 

- [How to create a streaming AsyncAPI microservice using golang](https://quobix.com/articles/asyncapi-streaming-using-golang/)

## First check out the repo

`git clone https://github.com/daveshanley/asyncapi-tutorials.git`

---

## Streaming Example Quickstart

1. Change directory to the streaming example.

`cd asyncapi-tutorials/streaming`

2. Boot the server and service.

`go run server/server.go`

3. Open a new terminal window and run the client and watch a stream of 10 random words appear, one every second

`go run client.go`
