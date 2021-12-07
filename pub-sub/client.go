package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/plank/services"
	"github.com/vmware/transport-go/plank/utils"
)

func main() {

	// create a message broker connector config and connect to
	// localhost over WebSocket on port 30080.
	config := &bridge.BrokerConnectorConfig{
		Username:   "guest",           // not required for demo, but our API requires it.
		Password:   "guest",           // ^^ same.
		ServerAddr: "localhost:30080", // our local plank instance, running RandomWordService
		UseWS:      true,              // connect over websockets
		WebSocketConfig: &bridge.WebSocketConfig{ // configure websocket
			WSPath: "/ws", // websocket endpoint
			UseTLS: false, // this isn't required locally, but you will always need to use it when connecting anywhere else.
		}}

	// get a pointer to transport
	b := bus.GetBus()

	// get a pointer to transport's channel manager
	cm := b.GetChannelManager()

	// connect to localhost:30080
	c, err := b.ConnectBroker(config)
	if err != nil {
		utils.Log.Fatalf("unable to connect to %s, error: %v", config.ServerAddr, err.Error())
	}

	// create local channels for pub-sub comms that are bridged to our joke-service channel.
	jokeSubChan := "jokes"
	cm.CreateChannel(jokeSubChan)

	// create a handler that will listen for a single response and then unsubscribe.
	jokeSubHandler, _ := b.ListenOnce(jokeSubChan)

	// mark our local 'jokes' channel as 'galactic' and map it to our connection and
	// the destinations defined by the AsyncAPI contract
	cm.MarkChannelAsGalactic(jokeSubChan, "/queue/joke-service", c)

	// create a wait group so our client stays running whilst we wait for a response.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start listening for our joke response.
	jokeSubHandler.Handle(
		func(msg *model.Message) {

			// extract our Joke response
			var joke services.Joke
			if err := msg.CastPayloadToType(&joke); err != nil {
				fmt.Printf("failed to cast payload: %s\n", err.Error())
			} else {
				// log out our joke to the console.
				utils.Log.Info(joke.Joke)
			}
			wg.Done()
		},
		func(err error) {
			utils.Log.Errorf("error received on channel: %e", err)
			wg.Done()
		})

	// create a joke request.
	req := &model.Request{Request: "get-joke"}
	reqBytes, _ := json.Marshal(req)

	// publish joke request
	c.SendJSONMessage("/pub/queue/joke-service", reqBytes)

	// wait for joke response to come in and be printed to the console.
	wg.Wait()

	// mark channels as local (unsubscribe)
	cm.MarkChannelAsLocal(jokeSubChan)

	// disconnect from our broker.
	c.Disconnect()
}
