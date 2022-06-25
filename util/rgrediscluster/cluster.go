package rgrediscluster



import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	Cluster *redis.ClusterClient
	ClusterCtx = context.Background()
)

type ClusterConfig struct{
	Host []string `yaml:"host"`
	Persistent bool `yaml:"persistent"`
	Auth string `yaml:"auth"`
	TimeOut time.Duration `yaml:"timeout"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
	Prefix string `yaml:"prefix"`
	PoolSize int `yaml:"poolSize"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}


func Setup() {




	Cluster = redis.NewClusterClient(&redis.ClusterOptions{

		//-------------------------------------------------------------------------------------------
		//集群相关的参数

		//集群节点地址，理论上只要填一个可用的节点客户端就可以自动获取到集群的所有节点信息。但是最好多填一些节点以增加容灾能力，因为只填一个节点的话，如果这个节点出现了异常情况，则Go应用程序在启动过程中无法获取到集群信息。


		//Addrs: setting.Conf.Cluster.Host,

		Addrs: []string{"192.168.1.85:6379","192.168.1.85:6380","192.168.1.86:6379","192.168.1.86:6380","192.168.1.87:6379","192.168.1.87:6380"},


		//Password: setting.Conf.Cluster.Auth,
		Password: "h0ScbctSdA",


		MaxRedirects: 8, // 当遇到网络错误或者MOVED/ASK重定向命令时，最多重试几次，默认8

		//只含读操作的命令的"节点选择策略"。默认都是false，即只能在主节点上执行。
		ReadOnly: false, // 置为true则允许在从节点上执行只含读操作的命令
		// 默认false。 置为true则ReadOnly自动置为true,表示在处理只读命令时，可以在一个slot对应的主节点和所有从节点中选取Ping()的响应时长最短的一个节点来读数据
		RouteByLatency: false,
		// 默认false。置为true则ReadOnly自动置为true,表示在处理只读命令时，可以在一个slot对应的主节点和所有从节点中随机挑选一个节点来读数据
		RouteRandomly: false,


		//------------------------------------------------------------------------------------------------------
		//ClusterClient管理着一组redis.Client,下面的参数和非集群模式下的redis.Options参数一致，但默认值有差别。
		//初始化时，ClusterClient会把下列参数传递给每一个redis.Client

		//每一个redis.Client的连接池容量及闲置连接数量，而不是clusterClient总体的连接池大小。实际上没有总的连接池
		//而是由各个redis.Client自行去实现和维护各自的连接池。
		//PoolSize:     setting.Conf.Cluster.PoolSize, // 连接池最大socket连接数，默认为5倍CPU数， 5 * runtime.NumCPU
		PoolSize:    150, // 连接池最大socket连接数，默认为5倍CPU数， 5 * runtime.NumCPU
		MinIdleConns: 10,                                 //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//超时
		DialTimeout:  5 * time.Second,                                     //连接建立超时时间，默认5秒。
		//ReadTimeout:  setting.Conf.Cluster.ReadTimeout * time.Second, //读超时，默认3秒， -1表示取消读超时
		//WriteTimeout: setting.Conf.Cluster.TimeOut * time.Second,     //写超时，默认等于读超时，-1表示取消读超时
		ReadTimeout: 2 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 2 * time.Second,     //写超时，默认等于读超时，-1表示取消读超时
		PoolTimeout:  4 * time.Second,                                     //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second,                                    //闲置连接检查的周期，无默认值，由ClusterClient统一对所管理的redis.Client进行闲置连接检查。初始化时传递-1给redis.Client表示redis.Client自己不用做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		//IdleTimeout:        setting.Conf.Cluster.IdleTimeout * time.Minute, //闲置超时，默认5分钟，-1表示取消闲置超时检查
		IdleTimeout:       5 * time.Minute, //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,                                     //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
	})
}
//
//func CreateKey(key string) string {
//	return setting.Conf.Cluster.Prefix + ":" + key
//}

