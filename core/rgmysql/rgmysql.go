package rgmysql

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"
	"github.com/jackylee92/rgo/core/rgglobal/rgerror"
	"github.com/jackylee92/rgo/core/rglog"

	logformat "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Factory struct {
	UniqId string
	logger interface{}
}

type client struct {
	Db     *gorm.DB
	Error  error
	dbName string
	logger interface{}
}

const defaultDBName = "default"

var connectPool map[string]*gorm.DB

type TableInterface interface {
	TableName() string
}

/*
* @Content : 开启
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-20
 */
func Start() {
	connect()
}

func New(uniqId string, logger interface{}) *Factory {
	if logger == nil {
		logger = rglog.New("nil")
	}
	return &Factory{UniqId: uniqId, logger: logger}
}

func (f *Factory) New(dbName string) (*client, error) {
	client := new(client)
	client.logger = f.logger
	return client.SetDB(dbName)
}

/*
* @Content : 链接
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-20
 */
func connect() {
	dsnList := getConfig()
	if len(dsnList) == 0 {
		return
	}
	logger := Logger{}
	switch rgconfig.GetStr(rgconst.ConfigMysqlLog) {
	case "debug":
		logger.LogLevel = 4
	case "info":
		logger.LogLevel = 4
	case "warn":
		logger.LogLevel = 3
	case "error":
		logger.LogLevel = 2
	default:
		logger.LogLevel = 2 // 默认只打印错误
	}

	logger.SlowThreshold = 200 * time.Millisecond
	formatter := &logformat.TextFormatter{
		// 不需要彩色日志
		DisableColors: true,
		// 定义时间戳格式

		TimestampFormat: rgconst.GoTimeFormat,
	}
	logformat.SetFormatter(formatter)
	connectPool = make(map[string]*gorm.DB, 5)
	connectNames := make([]string, 0, len(dsnList))
	for name, dsn := range dsnList {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})
		if err != nil {
			panic("数据库【" + name + "】链接失败")
		}
		newDB, err := db.DB()
		newDB.SetMaxOpenConns(int(rgconfig.GetInt(rgconst.ConfigMysqlMaxOpenConns)))
		newDB.SetMaxIdleConns(int(rgconfig.GetInt(rgconst.ConfigMysqlMaxIdleConns)))
		newDB.SetConnMaxLifetime(time.Duration(rgconfig.GetInt(rgconst.ConfigMysqlMaxLifetime)) * time.Second)
		connectPool[name] = db
		connectNames = append(connectNames, "【"+name+"】")
	}
	rglog.SystemInfo("启动项【mysql】", strings.Join(connectNames, ""), ":成功")
}

/*
* @Content : 获取根据不同配置文件方式数据链接配置
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-05-12
 */
func getConfig() (data map[string]string) {
	if rgconfig.Config.Config == "apollo" {
		dsnJson := rgconfig.GetStr(rgconst.ConfigKeyMysql)
		if dsnJson == "" {
			return data
		}
		err := json.Unmarshal([]byte(dsnJson), &data)
		if err != nil {
			rglog.SystemError("数据库配置文件解析失败")
			return data
		}
	}
	if rgconfig.Config.Config == "file" {
		data = rgconfig.GetStrMap(rgconst.ConfigKeyMysql)
	}
	return data
}

/*
* @Content : 设置数据库链接
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-23
 */
func (c *client) SetDB(param string) (*client, error) {
	if param == "" {
		param = defaultDBName
	}
	c.dbName = param
	conn, ok := connectPool[param]
	if !ok {
		c.Error = errors.New(rgerror.ErrorMysqlConnNil)
		return c, c.Error
	}
	ctx := context.WithValue(context.Background(), "logger", c.logger)
	c.Db = conn.WithContext(ctx)
	return c, nil
}
