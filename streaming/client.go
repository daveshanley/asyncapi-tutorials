package main

import (
	"sync"

	"github.com/vmware/transport-go/bridge"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/model"
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

	// create a local channel on the bus that we want to listen to in our application.
	myLocalChan := "my-local-word-stream"
	cm.CreateChannel(myLocalChan)

	// listen to stream of messages coming in on channel, a handler is returned
	// that allows you to add in lambdas that handle your success messages, and your errors.
	handler, _ := b.ListenStream(myLocalChan)

	// mark our local 'my-local-word-stream' myLocalChan as 'galactic' and map it to our connection and
	// the /topic/random-word service
	err = cm.MarkChannelAsGalactic(myLocalChan, "/topic/random-word", c)
	if err != nil {
		utils.Log.Fatalf("unable to map local channel to broker destination: %e", err)
	}

	// create a wait group that will wait 10 times before completing.
	var wg sync.WaitGroup
	wg.Add(10)

	// start and keep listening
	handler.Handle(
		func(msg *model.Message) {

			var randomWord string
			msg.CastPayloadToType(&randomWord)

			// log it out.
			utils.Log.Infof("Random word: %s", randomWord)

			wg.Done()
		},
		func(err error) {
			utils.Log.Errorf("error received on channel: %e", err)
		})

	// wait for 10 ticks of the stream, then we're done.
	wg.Wait()

	// close our handler, we're done.
	handler.Close()

	// mark channel as local (unsubscribe from /topic/random-word)
	cm.MarkChannelAsLocal(myLocalChan)

	// disconnect from our broker.
	c.Disconnect()
}
