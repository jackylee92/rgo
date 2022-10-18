package rgcoretest

import (
	"github.com/jackylee92/rgo"
	"testing"
)

//go test -v -run Test_Log core/rgcoretest/log_test.go -count=1 -config=../../config.yaml
func Test_Log(t *testing.T) {
	rgo.This.Log.Info(">>>>>>>>>>>>>>>>>>>>>>>>test")
}
