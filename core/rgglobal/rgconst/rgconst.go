package rgconst

const (
	ProcessKilled string = "收到信号，进程被结束"
	// BaseConfigEnd string = "/config.json"
	UnknownError string = "未知错误"
	// StartFuncMaxCount <LiJunDong : 2022-05-30 18:10:42> --- 注册的启动函数
	StartFuncMaxCount = 20

	// ReturnSuccessCode <LiJunDong : 2022-03-10 19:03:36> --- 请求常用返回状态
	ReturnSuccessCode int64  = 200
	ReturnPanicCode   int64  = 500
	GoTimeFormat      string = "2006-01-02 15:04:05" // go 完整时间格式
	GoDateFormat      string = "2006-01-02"          // go 日期格式

	// UniqIDKey <LiJunDong : 2022-03-28 11:24:00> --- context中uniqid的key名
	ContextUniqIDKey        string = "uniqid"
	ContextStartTimeKey     string = "start_time"
	ContextContainerKey     string = "container"
	ContextJeargerCtxKey    string = "jearger_ctx"
	ContextJeargerTracerKey string = "jearger_tracer"

	//  <LiJunDong : 2022-06-18 14:11:50> --- 心跳的url
	ConfigHeartBeatUrl string = "elb-status"
	//  <LiJunDong : 2022-06-18 14:11:50> --- 设置日志级别
	ConfigSetLogLevelUrl string = "setLogLevel"
	ConfigGetLogLevelUrl string = "getLogLevel"
	ConfigGetUrl         string = "getConfig"

	// <LiJunDong : 2022-03-28 11:24:00> --- 系统配置文件名
	ConfigKeyLogLevel             string = "sys_log_level"               // 日志级别
	ConfigKeyMysql                string = "sys_mysql"                   // mysql
	ConfigKeyAppName              string = "sys_app_name"                // 项目名称
	ConfigKeyPort                 string = "sys_port"                    // 端口
	ConfigKeyMessage              string = "sys_language"                // 语言
	ConfigKeyJaergerStatus        string = "sys_jaerger_status"          // jaerger开关
	ConfigKeyJaergerHost          string = "sys_jaerger_host"            // jaerger配置地址
	ConfigKeyRequestLog           string = "sys_request_log"             // 请求日志
	ConfigKeyDebug                string = "sys_debug"                   // 调试模式，开启调试模式，所有的日志将会在终端打印
	ConfigKeyHttpAllowCrossDomain string = "sys_http_allow_cross_domain" // 是否允许跨域
	ConfigKeyRedis                string = "sys_redis"                   // apollo redis配置
	ConfigKeyRedisType            string = "sys_redis_type"              // 本地redis配置
	ConfigKeyRedisAddr            string = "sys_redis_addr"              // 本地redis 地址
	ConfigKeyRedisAddrs           string = "sys_redis_addrs"             // 本地redis 集群地址
	ConfigKeyRedisPassword        string = "sys_redis_password"          // 本地redis认证
	ConfigKeyRedisDB              string = "sys_redis_db"                // 本地redisDB
	ConfigKeyRedisPoolSize        string = "sys_redis_pool_size"         // 本地redis连接池最大socket连接数，默认为10倍CPU数， 10 * runtime.NumCPU
	ConfigKeyRedisMinIdleConns    string = "sys_redis_min_idle_conns"    // 本地redis在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
	ConfigReturnUniqId            string = "sys_return_uniqid"           // 是否返回uniqid
	ConfigMysqlLog                string = "sys_mysql_log"               // mysql日志级别
	ConfigMysqlMaxOpenConns       string = "sys_mysql_max_open_conns"    // mysql初始化连接数
	ConfigMysqlMaxIdleConns       string = "sys_mysql_max_idle_conns"    // mysql最大连接数
	ConfigMysqlMaxLifetime        string = "sys_mysql_max_life_time"     // mysql链接生命时长
	ConfigLogDir                  string = "sys_log_dir"                 // 日志目录
	ConfigPProf                   string = "sys_pprof"                   // pprof开启关闭
)
