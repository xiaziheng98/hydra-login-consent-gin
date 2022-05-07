package hydra

import client "github.com/ory/hydra-client-go"

var Hydra *client.APIClient

func InitHydra() {
	configuration := client.NewConfiguration()
	configuration.Servers = []client.ServerConfiguration{
		{
			URL: "http://localhost:4445", // Admin API URL
		},
	}
	Hydra = client.NewAPIClient(configuration)
}
