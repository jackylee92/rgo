package rgerror

const (
	ErrorBasePath          string = "初始化项目根目录失败"
	ErrorStartConfigExists string = "-config配置文件不能为空"
	ErrorConfigExist       string = "初始化失败,配置文件未读取到"
	ErrorConfigTypeExist   string = "初始化失败,配置文件类型错误:file/apollo"
	ErrorConfigParse       string = "配置文件解析错误"
	ErrorAppConfig         string = "项目配置文件错误"
	ErrorBaseConfig        string = "框架配置文件读取错误"
	ErrorConfigInit        string = "项目配置启动失败"
	ErrorRequirePanic      string = "程序异常panic"
	ErrorAppNameNil        string = "未获取到AppName"
	ErrorJaegerHostNil     string = "未获取到JaegerHost"
	ErrorLocalIpNil        string = "获取本地IP失败"
	ErrorStartFuncOverflow string = "启动函数过多"

	// model相关
	ErrorMysqlConnNil   string = "未配置该数据库链接"
	ErrorRedisConfigNil string = "未获取到Redis链接配置"
	ErrorRedisClientNil string = "获取redis链接对象失败"

	// curl相关
	CurlErrorServerSelf   string = "程序内部异常"
	CurlErrorServerSelfEn string = "Program internal error"
	CurlSuccess           string = "请求成功"
	CurlSuccessEn         string = "Success"
	Curl404Error          string = "地址错误404,^-^"
)
