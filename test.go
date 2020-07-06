package desultory

import (
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"
)

var offlineEnvironmentVariable = "OFFLINE"

type testStruct struct {
	Text string
	Number int
}

func getOffline() bool {
	ev := os.Getenv(offlineEnvironmentVariable)
	if strings.EqualFold(ev, "true") {
		return true
	}
	return false
}

func getTestDirectory() (string, error) {
	tid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	tp := path.Join("/tmp/" + tid.String())
	e, err := DirectoryOrFileExists(tp)
	if err != nil {
		return "", err
	}
	if e {
		err = os.RemoveAll(tp)
		if err != nil {
			return "", err
		}
	}
	err = os.MkdirAll(tp, 0777)
	if err != nil {
		return "", err
	}
	return tp, nil
}

var testStartTime time.Time

func setupTest(t *testing.T) {
	testStartTime = time.Now()
	v := reflect.ValueOf(t)
	tn := v.FieldByName("name")
	logrus.Infof("started test '%v' at '%v'", tn, testStartTime)
}

func tearDownTest(t *testing.T) {
	ts := time.Now().Sub(testStartTime)
	v := reflect.ValueOf(t)
	tn := v.FieldByName("name")
	logrus.Infof("completed test '%v' after '%v' seconds", tn, ts)
}