package services

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/service"
)

const (
	JokeServiceChannel = "joke-service" // matches asyncapi destination channel.
)

// Joke is a representation of what is returned by our providing JokeAPI.
type Joke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// JokeService will return a terrible joke, on demand.
type JokeService struct{}

// NewJokeService will return an instance of JokeService.
func NewJokeService() *JokeService {
	return &JokeService{}
}

// Init will fire when the service is being registered by Plank.
func (js *JokeService) Init(core service.FabricServiceCore) error {

	// JokeService always returns JSON objects as responses. Set default 'application/json' content-type headers.
	core.SetDefaultJSONHeaders()
	return nil
}

// HandleServiceRequest will handle icoming requests from event bus on our service channel.
func (js *JokeService) HandleServiceRequest(request *model.Request, core service.FabricServiceCore) {
	switch request.Request {
	case "get-joke":
		js.getJoke(request, core)
	default:
		core.HandleUnknownRequest(request)
	}
}

// getJoke calls our terrible joke service, and returns the response or error back to the requester.
func (js *JokeService) getJoke(request *model.Request, core service.FabricServiceCore) {

	// make API call using inbuilt RestService to make network calls. Use the wonderful https://icanhazdadjoke.com API.
	core.RestServiceRequest(&service.RestServiceRequest{
		Uri:    "https://icanhazdadjoke.com",
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		ResponseType: reflect.TypeOf(&Joke{}),
	}, func(response *model.Response) {

		// send back a successful joke.
		core.SendResponse(request, response.Payload.(*Joke))

	}, func(response *model.Response) {

		// something went wrong with the API call, tell the requester.
		fabricError := service.GetFabricError("Get Joke API Call Failed", response.ErrorCode, response.ErrorMessage)
		core.SendErrorResponseWithPayload(request, response.ErrorCode, response.ErrorMessage, fabricError)
	})
}

// GetRESTBridgeConfig returns a config for a REST endpoint for this Joke Service
func (js *JokeService) GetRESTBridgeConfig() []*service.RESTBridgeConfig {
	return []*service.RESTBridgeConfig{
		{
			ServiceChannel: JokeServiceChannel, // where is this service running?
			Uri:            "/rest/joke",       // what path do you want to map to this service?
			Method:         http.MethodGet,     // which method on this path should we map?
			AllowHead:      true,               // can the client send a HEAD request on this path?
			AllowOptions:   true,               // can the client send an OPTIONS request on this path?
			FabricRequestBuilder: func(w http.ResponseWriter, r *http.Request) model.Request {

				// Plank will essentially call this service like any other for every inbound HTTP request
				// so we create a message on behalf of the client.
				return model.Request{
					Id:                &uuid.UUID{},
					Request:           "get-joke", // which command do we want to run?
					BrokerDestination: nil,        // don't worry anout this, unless using muliple brokers in your config.
				}
			},
		},
	}
}

// OnServerShutdown is not implemented in this service.
func (js *JokeService) OnServerShutdown() {}

// OnServiceReady has no functionality in this service, so it returns a pre-fired channel.
func (js *JokeService) OnServiceReady() chan bool {

	// ready right away, nothing to do.
	readyChan := make(chan bool, 1)
	readyChan <- true
	return readyChan
}
