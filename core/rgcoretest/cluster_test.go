package rgcoretest

import (
	"encoding/json"
	"fmt"
	"github.com/jackylee92/rgo/core/rgmodel/rgcluster"
	"github.com/jackylee92/rgo/library/rguser"
	"github.com/magiconair/properties/assert"
	"testing"
)

//go test -v -run Test_Log core/rgcoretest/cluster_test.go -count=1 -config=../../config.yaml
func TestCluster(t *testing.T) {
	config := rgcluster.ClusterConfig{
		Host:        []string{"192.168.1.85:6379", "192.168.1.85:6380", "192.168.1.86:6379", "192.168.1.86:6380", "192.168.1.87:6379", "192.168.1.87:6380"},
		Auth:        "h0ScbctSdA",
		TimeOut:     20,
		ReadTimeout: 10,
		Prefix:      "CustomerCenterTest",
		PoolSize:    20,
		IdleTimeout: 5,
	}
	var memberInfo rguser.User

	rgcluster.Setup(config)

	key := rgcluster.CreateKey("user_base:11819")

	fmt.Println("key", key)
	res, err := rgcluster.Cluster.Get(rgcluster.ClusterCtx, key).Result()

	assert.Equal(t, err == nil, true, "用户信息不存在")

	err1 := json.Unmarshal([]byte(res), &memberInfo)
	assert.Equal(t, err1 == nil, true, "用户信息解析失败")
	fmt.Println("memberInfo", memberInfo)
}
