package rgcoretest

import (
	"fmt"
	_ "github.com/jackylee92/rgo"
	"github.com/jackylee92/rgo/core/rgredis"
	"testing"
)

//go test -v -run TestRedisOne core/rgcoretest/redis_test.go -count=1 -config=../../config.yaml
func TestRedisOne(t *testing.T) {
	config := rgredis.Config{
		Addr:     "test.redis.ruigushop.com:6379",
		Password: "h0ScbctSdA",
		TimeOut:  20,
		PoolSize: 20,
	}
	redisPool, err := rgredis.Start(config)
	if err != nil {
		fmt.Println("连接失败|" + err.Error())
		return
	}
	redisClient, _ := rgredis.Pool(redisPool)
	if err != nil {
		fmt.Println("获取链接失败|" + err.Error())
		return
	}
	ok, err := redisClient.SIsMember("222", "123")
	if err != nil {
		fmt.Println("error", err.Error())
	}else {
		fmt.Println(ok)
	}
	return



	//config2 := rgredis.Config{
	//	Addr:     "test.redis.ruigushop.com:6379",
	//	Password: "h0ScbctSdA",
	//	TimeOut:  20,
	//	PoolSize: 20,
	//}
	//redisPool2, err := rgredis.Start(config2)
	//if err != nil {
	//	fmt.Println("连接失败|" + err.Error())
	//	return
	//}
	//redisClient2, _ := rgredis.Pool(redisPool2)
	//if err != nil {
	//	fmt.Println("获取链接失败|" + err.Error())
	//	return
	//}
	//data, err := redisClient2.GetClient()
	//fmt.Println(data, err)
}
