package rgcoretest

import (
	"fmt"
	_ "github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/core/rgrediscluster"
	"testing"
)

//go test -v -run TestRedisCluster core/rgcoretest/redis_cluster_test.go -count=1 -config=../../config.yaml
func TestRedisCluster(t *testing.T) {
	config := rgrediscluster.Config{
		Addrs:    []string{"192.168.1.85:6379", "192.168.1.85:6380", "192.168.1.86:6379", "192.168.1.86:6380", "192.168.1.87:6379", "192.168.1.87:6380"},
		Password: "h0ScbctSdA",
		TimeOut:  20,
		PoolSize: 20,
	}
	redisPool, err := rgrediscluster.Start(config)
	if err != nil {
		fmt.Println("连接失败|" + err.Error())
		return
	}
	redisClient, _ := rgrediscluster.Pool(redisPool)
	if err != nil {
		fmt.Println("获取链接失败|" + err.Error())
		return
	}
	ok, err := redisClient.Setex("lilili", "123", 0)
	fmt.Println(ok, err)
}
