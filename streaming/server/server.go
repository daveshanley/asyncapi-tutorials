package main

import (
	"os"

	"github.com/daveshanley/asyncapi-tutorials/tutorial-1/services"
	"github.com/vmware/transport-go/plank/pkg/server"
	"github.com/vmware/transport-go/plank/utils"
)

// main will create a new instance of plank using a default configuration.
func main() {

	// create a default server configuration.
	serverConfig, err := server.CreateServerConfig()
	if err != nil {
		utils.Log.Fatalln(err)
		return
	}

	// create a new platform server from our configuration.
	platformServer := server.NewPlatformServer(serverConfig)

	// register our RandomWordService with our platform server.
	if err = platformServer.RegisterService(services.NewRandomWordService(), services.RandomWordChannel); err != nil {
		utils.Log.Fatalln(err)
		return
	}

	// register a system channel with the platform, so we can catch interrupts and shut down cleanly.
	syschan := make(chan os.Signal, 1)

	// start plank and start streaming random words to everyone.
	platformServer.StartServer(syschan)
}
