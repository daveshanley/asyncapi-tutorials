# AsyncAPI Tutorials for Golang

This code accompanies articles on [quobix.com](https://quobix.com). If you would like to gain full context behind this code and learn more about the tooling it uses, then check out the following articles. 

- [How to create an event-driven API via AsyncAPI using golang](https://quobix.com/articles/asyncapi-pubsub-using-golang/)
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

> make sure you stop the server process before trying other examples.

---

## Pub-Sub Example Quickstart


1. Change directory to the pub-sub example.

`cd asyncapi-tutorials/pub-sub`

2. Boot the server and service.

`go run server/server.go &`

3. Open a new terminal window and run the client. Enjoy your terrible joke.

`go run client.go`

