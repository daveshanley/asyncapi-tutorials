package main

import (
	"encoding/json"
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
	jokeSubChan := "joke-service"
	cm.CreateChannel(jokeSubChan)

	// create a handler that will listen for a single response and then unsubscribe.
	jokeSubHandler, _ := b.ListenOnce(jokeSubChan)

	// mark our local 'joke-sub' and 'joke-pub' as 'galactic' and map it to our connection and
	// the destinations defined by the AsyncAPI contract
	//cm.MarkChannelAsGalactic(jokePubChan, "/pub/queue/joke-service", c)
	cm.MarkChannelAsGalactic(jokeSubChan, "/queue/joke-service", c)

	// create a wait group so our client stays running whilst we wait for a response.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start listening for our joke response.
	jokeSubHandler.Handle(
		func(msg *model.Message) {

			// unmarshal the message payload into a model.Response object
			// this is a wrapper transport uses when being used as a server,
			// it encapsulates a rich set of data about the message,
			// but you only really care about the payload (body)
			r := &model.Response{}
			d := msg.Payload.([]byte)
			err := json.Unmarshal(d, &r)
			if err != nil {
				utils.Log.Errorf("error unmarshalling request: %v", err.Error())
				return
			}

			// the value we want is in the payload of our model.Response
			value := r.Payload.(services.Joke)

			// log out our joke to the console.
			utils.Log.Info(value.Joke)

			wg.Done()
		},
		func(err error) {
			utils.Log.Errorf("error received on channel: %e", err)
			wg.Done()
		})

	// publish joke request
	c.SendJSONMessage("/pub/queue/joke-service", []byte("pizza"))

	// wait for joke response to come in and be printed to the console.
	wg.Wait()

	// mark channels as local (unsubscribe)
	cm.MarkChannelAsLocal(jokeSubChan)

	// disconnect from our broker.
	c.Disconnect()
}
