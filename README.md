# RGO

## RGo：

Rgo是一个基于Gin二次开发的框架，该框架集成了Log、MySQL、Jaeger、Config、pprof、Redis、Hook等组件功能，所有组件配置简单，高并发场景下日志文件写入，MySQL单例模式等。

## 安装

``go get github.com/jackylee92/rgo``

#### 建议目录结构

```
├── README.md              // 项目ReadMe
├── cmd                    // 项目主入口
│   └── main.go 
├── common                 // 项目公共方法/常量
│   ├── README.md
│   ├── common.go
│   └── const.go
├── config                 // 配置文件
│   └── config.yaml
├── controller             // 控制器层
│   ├── README.md
│   └── api.go
├── docs                   // 详细的迭代等文档
│   └── 项目介绍.md
├── go.mod
├── go.sum
├── hook                   // 钩子文件
│   └── hook.go
├── message                // 语言包定义文件
│   ├── README.md
│   └── message.go
├── middleware             // 中间件
│   └── README.md
├── model                  // 数据层
│   ├── README.md
│   ├── datamate           // 调用三方数据源层
│   │   └── datemate.go
│   └── mysql              // 对接MySQL层
│       ├── README.md
│       ├── userdb
│       │   └── user.go
│       └── mysql.go
├── router                 // 路由层
│   ├── README.md
│   └── router.go
├── service                // 业务层
│   ├── README.md
│   └── user.go
├── storage                // 静态文件/日志
│   ├── README.md
│   └── log
│       └── log.log
├── test                   // 单元测试
│   ├── README.md
│   └── test_user.go
└── validator              // 验证层
    └── validator.go
```

#### 启动：

 ``go run cmd/main.go -config=config/config.yaml``

## 模块使用

### Router

对Gin路由二次封装，``router := rgrouter.NewRouter()``获取路由对象。NewRouter自动加载``eager.Handle(), requestlog.Handle(), recovery.Handle(), container.Handle()`` 中间件，定义默认路由：健康检查``elb-status``,动态修改日志级别``setLogLevel``等，可查看源码``core/rgrouter/rgrouter.go``

案例：

```go
// 获取默认路由
router := rgrouter.NewRouter()
// 定义路由规则
api := router.Group("/v1")
{
   api.GET("info/get", info.GetHandle) 
}
// 启动
rgrouter.Run(router)
```

### Log

日志默认会记录在storage/log/项目名.log，默认会按天备份，项目名定义在config.yaml中``sys_app_name``，日志使用``zerolog``默认每秒每5000条输入一次。

日志输入流只有一个，在12点时会重命名，建立新的项目名.log文件，更新输入目标为新的日志文件。

通过获取请求上下文``rgrequest.Get(ctx)``写入日志会记录本次唯一的uniqId。

配置日志：``config/config.yaml``

日志文件配置：sys_log_dir：xxx/xxx/

日志级别配置：sys_log_level:debug/info/warn/error

#### 记录日志

```go
this := rgrequest.Get(ctx) // 获取请求上下文
this.Log.Info("info") // 记录info日志
this.Log.Error(err) // 记录错误日志
```

### MySQL

配置数据库连接 config/config.yaml

```yaml
sys_mysql:
   db_name:"root:password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
```

获取MySQL链接

```go
model, err := this.Mysql.New("db_name") // 获取db_name链接
if err != nil {
    return err
}
```

操作

```go
 92 func (c *Coupon) Update(param mysql.SearchParam, data map[string]interface{}) (err error) {
 93     model, err := param.This.Mysql.New("db_name")
 95     if err != nil {
 96         return err
 97     }
 98     data["updated_at"] = mysql.NowTime()
 99     return model.Db.Table(c.TableName()).Where(param.Query, param.Args...).Updates(data).Error
100 }
101
102 func (c *Coupon) Find(param mysql.SearchParam) (err error) {
103     model, err := param.This.Mysql.New("db_name")
104     if err != nil {
105         return err
106     }
108     m := model.Db.Table(c.TableName()).Where(param.Query, param.Args...)
109     if param.Fields != nil && len(param.Fields) != 0 {
110         m = m.Select(param.Fields)
111     }
112     m.Find(&c)
113     if m.RowsAffected == 0 {
114         return mysql.ErrNil
115     }
116     return m.Error
117 }
```

### Config

日志文件实现了本地yaml/apollo，不同的配置方式，只要实现ConfigInterface即可，通过配置config.yaml中 ``config: file``表示配置使用文件方式

```go
  8 type ConfigInterface interface {
  9     GetStr(string) string               // 获取String类型
 10     GetInt(string) int64                // 获取int64类型
 11     GetBool(string) bool                // 获取bool类型
 12     GetStrMap(string) map[string]string // 获取map[string]string类型
 13     GetStrSlice(string) []string        // 获取[]string类型
 14     GetContent() string                 // 获取所有配置文件string
 15     SetConfig(string, interface{})      // 设置配置
 16     Isset(string) bool                  // 判断配置是否存在
 17     Reload() error                      // 重新加载配置
 18     Load() func() error
 19 }
```

##### 使用

```go
name := rgconfig.GetStr("name")
age := rgconfig.GetInt("age")
res := rgconfig.GetBool("res")
```

### pprof

通过配置文件，是否启动 ``sys_pprof`` true/false

### Redis

### Hook

启动时可像rgstarthook中注册函数

```go
import (
   "github.com/jackylee92/rgo/core/rgstarthook"
    _ "github.com/jackylee92/rgo/util/rgstarthook" // 注册内置启动通知
) 
rgstarthook.Run() // 运行hook注册的方法
```