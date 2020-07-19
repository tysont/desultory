package desultory

import (
	f "github.com/fauna/faunadb-go/faunadb"
)

func CreateFaunaDatabase(database string) error {
	if faunaRootClient == nil {
		initializeFauna()
	}
	_, err := faunaRootClient.Query(
		f.If(
			f.Exists(f.Database(database)),
			true,
			f.CreateDatabase(f.Obj{"name": database})))
	if err != nil {
		return err
	}
	res, err := faunaRootClient.Query(
		f.CreateKey(
			f.Obj{"database": f.Database(database), "role": "server"}))
	if err != nil {
		return err
	}
	var key string
	res.At(f.ObjKey("secret")).Get(&key)
	client := faunaRootClient.NewSessionClient(key)
	faunaDatabaseClients[database] = client
	return nil
}

func DeleteFaunaDatabase(database string) error {
	if faunaRootClient == nil {
		initializeFauna()
	}
	_, err := faunaRootClient.Query(
		f.If(
			f.Exists(f.Database(database)),
			f.Delete(f.Database(database)),
			false))
	if err != nil {
		return err
	}
	faunaDatabaseClients[database] = nil
	return nil
}

func CreateFaunaCollection(database string, collection string) error {
	client := faunaDatabaseClients[database]
	_, err := client.Query(
			f.If(
				f.Exists(f.Collection(collection)),
				true,
				f.CreateCollection(f.Obj{"name": collection})))
	return err
}

func DeleteFaunaCollection(database string, collection string) error {
	client := faunaDatabaseClients[database]
	_, err := client.Query(
		f.If(
			f.Exists(f.Collection(collection)),
			f.Delete(f.Collection(collection)),
			false))
	return err
}

func CreateFaunaIndex(database string, collection string, index string, pkey string) error {
	client := faunaDatabaseClients[database]
	_, err := client.Query(
		f.If(
			f.Exists(f.Index(index)),
			true,
			f.CreateIndex(f.Obj{
				"name": index,
				"source": f.Collection(collection),
				"unique": true,
				"terms": f.Arr{f.Obj{"field": f.Arr{"data", pkey}}}})))
	return err
}

func DeleteFaunaIndex(database string, index string) error {
	client := faunaDatabaseClients[database]
	_, err := client.Query(
		f.If(
			f.Exists(f.Index(index)),
			f.Delete(f.Index(index)),
			false))
	return err
}

func CreateFaunaInstance(database string, collection string, o interface{}) error {
	client := faunaDatabaseClients[database]
	_, err := client.Query(
		f.Create(
			f.Collection(collection),
			f.Obj{"data": o}))
	return err
}

func GetFaunaInstance(database string, index string, pkval string, o interface{}) error {
	client := faunaDatabaseClients[database]
	res, err := client.Query(
		f.Get(
			f.MatchTerm(f.Index(index), pkval)))
	if err != nil {
		return err
	}
	return res.At(f.ObjKey("data")).Get(o)
}