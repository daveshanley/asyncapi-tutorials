package services

import (
	"math/rand"
	"reflect"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/service"
)

const (
	RandomWordChannel = "random-word" // matches asyncapi destination channel.
)

// RandomWordService will broadcast a random word on the "simple-stream" channel, every one second.
type RandomWordService struct {
	words         []string                  // list of random words.
	transportCore service.FabricServiceCore // reference to transport services we will need later on.
	readyChan     chan bool                 // once we're ready, let plank know via this channel.
	cronJob       *cron.Cron                // cronjob that runs every 1s.
}

// NewRandomWordService will return a new instance of RandomWordService
func NewRandomWordService() *RandomWordService {
	return &RandomWordService{}
}

// Init will fire when our service is being registered by Plank. A pointer to service.FabricServiceCore is passed in
// so we can capture that pointer.
func (rws *RandomWordService) Init(core service.FabricServiceCore) error {

	// capture a reference to transport core services.
	rws.transportCore = core
	return nil
}

// OnServiceReady fires once Plank has all services loaded and ready to run. The method returns a bool channel that Plank
// will wait for a signal on and then make the service available.
func (rws *RandomWordService) OnServiceReady() chan bool {
	rws.readyChan = make(chan bool, 1)

	// fetch a list of random words (which runs asynchronously), so it immediately returns.
	rws.fetchRandomWords()

	return rws.readyChan
}

// fetchRandomWords will call a public REST endpoint that very kindly returns random words.
func (rws *RandomWordService) fetchRandomWords() {

	restRequest := &service.RestServiceRequest{
		Uri:          "https://random-word-api.herokuapp.com/word?number=500",
		Method:       "GET",
		ResponseType: reflect.TypeOf(rws.words),
	}

	// Transport provides a REST Service that makes this API call and provides handlers for the result.
	rws.transportCore.RestServiceRequest(restRequest,
		rws.handleWordFetchSuccess, // handle a successful API call.
		rws.handleWordFetchFailure) // handle a failed API call.
}

// handleWordFetchSuccess will parse a successful incoming word response from our source API.
func (rws *RandomWordService) handleWordFetchSuccess(response *model.Response) {

	// set the word list to the response returned by the REST API Call.
	rws.words = response.Payload.([]string)

	// start random word cron job.
	rws.fireRandomWords()

	// send a signal down our ready channel, so Plank knows to continue.
	rws.readyChan <- true
}

// handleWordFetchFailure will parse a failed random word API request.
func (rws *RandomWordService) handleWordFetchFailure(response *model.Response) {

	// now we have no data, so make something up using some hard coded values.
	rws.words = []string{"magnum", "fox", "kitty", "cotton", "ember"}

	// start random word cron job.
	rws.fireRandomWords()

	// we have a back up data-set loaded, so send a signal down our ready channel, so Plank knows to continue.
	rws.readyChan <- true
}

// fireRandomWords will create a cron job that repeats every minute, that sends a message to all subscribers
// every minute. We then capture a pointer to that cronjob on our RandomWordService.
func (rws *RandomWordService) fireRandomWords() {

	// function to fire every second.
	var fireMessage = func() {
		id := uuid.New()

		// send a message containing a random word.
		rws.transportCore.SendResponse(&model.Request{Id: &id}, rws.getRandomWord())
	}
	rws.cronJob = cron.New()
	rws.cronJob.AddFunc("@every 1s", fireMessage)
	rws.cronJob.Start()
}

// getRandomWord will return a random word from our in memory list.
func (rws *RandomWordService) getRandomWord() string {
	return rws.words[rand.Intn(len(rws.words)-1)]
}

// OnServerShutdown will stop the cronjob firing cleanly when Plank shuts down.
func (rws *RandomWordService) OnServerShutdown() {
	rws.cronJob.Stop()
}

// GetRESTBridgeConfig is not used by this service.
func (rws *RandomWordService) GetRESTBridgeConfig() []*service.RESTBridgeConfig {
	return nil
}

// HandleServiceRequest is not used by this servuce.
func (rws *RandomWordService) HandleServiceRequest(r *model.Request, c service.FabricServiceCore) {
	// do nothing in here, we're not listening for any requests.
}
