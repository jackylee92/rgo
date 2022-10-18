package rgtest

import (
	"github.com/jackylee92/rgo"
	"testing"
)

// go test -v -run TestStart test/rg_test.go -count=1 -config=../config.yaml
func TestStart(t *testing.T) {
	rgo.This.Log.Info("test", "start")
}
