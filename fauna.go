package desultory

import f "github.com/fauna/faunadb-go/faunadb"

var faunaRootAccessKey = "fnADxAKKZeACDXgY__82gfYHawD6S3sAiyeassqT"
var faunaRootClient *f.FaunaClient
var faunaDatabaseClients map[string]*f.FaunaClient

func initializeFauna() {
	faunaRootClient = f.NewFaunaClient(faunaRootAccessKey)
	faunaDatabaseClients = make(map[string]*f.FaunaClient, 0)
}
