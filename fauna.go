package desultory

import (
	f "github.com/fauna/faunadb-go/faunadb"
	"os"
)

var faunaRootAccessKeyEnvironmentVariable = "FAUNA_ACCESS_KEY"
var faunaRootClient *f.FaunaClient
var faunaDatabaseKeyReferences map[string]*f.RefV
var faunaDatabaseClients map[string]*f.FaunaClient

func initializeFauna() {
	k := getFaunaRootAccessKey()
	faunaRootClient = f.NewFaunaClient(k)
	faunaDatabaseKeyReferences = make(map[string]*f.RefV, 0)
	faunaDatabaseClients = make(map[string]*f.FaunaClient, 0)
}

func getFaunaRootAccessKey() string {
	k := os.Getenv(faunaRootAccessKeyEnvironmentVariable)
	return k
}