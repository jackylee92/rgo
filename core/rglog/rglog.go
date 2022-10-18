package rglog

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackylee92/rgo/core/rgconfig"
	"github.com/jackylee92/rgo/core/rgglobal"
	"github.com/jackylee92/rgo/core/rgglobal/rgconst"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var logLevel string

func Start() {
	zerolog.SetGlobalLevel(setLogLever())
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	if logLevel == "debug" {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}
	configLogDir := rgconfig.GetStr(rgconst.ConfigLogDir)
	if configLogDir == "" {
		logDir = filePathMerge(rgglobal.BasePath, "/storage/log/")
	} else {
		logDir = filePathMerge(configLogDir)
	}
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logDir, 0777); err != nil {
			panic("日志文件创建失败|err" + err.Error())
		}
	}
}

type Client struct {
	UniqId string
	sync.Pool
}

var logDir string

var levelSimpler *zerolog.LevelSampler = &zerolog.LevelSampler{
	DebugSampler: &zerolog.BurstSampler{
		Burst:       5000,
		Period:      time.Second,
		NextSampler: &zerolog.BasicSampler{N: 100},
	},
}

// New 获取一个新的对象
// @Param   : uniqid string
// @Return  : *Client
// @Author  : LiJunDong
// @Time    : 2022-06-03
func New(uniqId string) *Client {
	// TODO <LiJunDong : 2022-06-03 11:11:07> --- 使用pool避免频繁创建对象
	return &Client{UniqId: uniqId}
}

/*
* @Content : 日志记录
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-01
 */
func (c *Client) Info(any ...interface{}) {
	param := localDebug("Info", c.UniqId, any)
	if e := log.Info(); e.Enabled() {
		if param == "" {
			param = reqToStr(any)
		}
		nowDate := time.Now().Format(rgconst.GoDateFormat)
		filePath := filePathMerge(logDir, "/"+nowDate, "_INFO.log")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		output := zerolog.ConsoleWriter{
			Out:     f,
			NoColor: true,
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Local().Format(rgconst.GoTimeFormat)
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
		}
		logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
		logger.Info().Fields(map[string]interface{}{"UniqId": c.UniqId}).Msg(param)
	}
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func (c *Client) Error(any ...interface{}) {
	param := localDebug("Error", c.UniqId, any)
	if e := log.Error(); e.Enabled() {
		if param == "" {
			param = reqToStr(any)
		}
		nowDate := time.Now().Format(rgconst.GoDateFormat)
		filePath := filePathMerge(logDir, "/"+nowDate, "_ERROR.log")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		output := zerolog.ConsoleWriter{
			Out:     f,
			NoColor: true,
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Local().Format(rgconst.GoTimeFormat)
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
		}
		logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
		logger.Error().Fields(map[string]interface{}{"UniqId": c.UniqId}).Msg(param)
	}
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func (c *Client) Debug(any ...interface{}) {
	param := localDebug("Debug", c.UniqId, any)
	if e := log.Debug(); e.Enabled() {
		if param == "" {
			param = reqToStr(any)
		}
		nowDate := time.Now().Format(rgconst.GoDateFormat)
		filePath := filePathMerge(logDir, "/"+nowDate, "_INFO.log")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		output := zerolog.ConsoleWriter{
			Out:     f,
			NoColor: true,
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Local().Format(rgconst.GoTimeFormat)
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
		}
		logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
		logger.Debug().Msg(param)
	}
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func (c *Client) Warn(any ...interface{}) {
	param := localDebug("Warn", c.UniqId, any)
	if e := log.Warn(); e.Enabled() {
		if param == "" {
			param = reqToStr(any)
		}
		nowDate := time.Now().Format(rgconst.GoDateFormat)
		filePath := filePathMerge(logDir, "/"+nowDate, "_INFO.log")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		output := zerolog.ConsoleWriter{
			Out:     f,
			NoColor: true,
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Local().Format(rgconst.GoTimeFormat)
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
		}
		logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
		logger.Warn().Fields(map[string]interface{}{"UniqId": c.UniqId}).Msg(param)
	}
}

/*
* @Content :
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func (c *Client) Fatal(any ...interface{}) {
	param := localDebug("Fatal", c.UniqId, any)
	if e := log.Fatal(); e.Enabled() {
		if param == "" {
			param = reqToStr(any)
		}
		nowDate := time.Now().Format(rgconst.GoDateFormat)
		filePath := filePathMerge(logDir, "/"+nowDate, "_ERROR.log")
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			return
		}
		output := zerolog.ConsoleWriter{
			Out:     f,
			NoColor: true,
			FormatTimestamp: func(i interface{}) string {
				return time.Now().Local().Format(rgconst.GoTimeFormat)
			},
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatCaller: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s|", i)
			},
		}
		logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
		logger.Fatal().Fields(map[string]interface{}{"UniqId": c.UniqId}).Msg(param)
	}
}

/*
* @Content : Print
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-09
 */
func Print(param ...interface{}) {
	log.Print(param...)
}
func Println(param ...interface{}) {
	log.Print(param...)
}

/*
* @Content : 获取日志级别
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-10
 */
func setLogLever() (level zerolog.Level) {
	level = zerolog.InfoLevel
	logLevel = "info"
	switch rgconfig.GetStr(rgconst.ConfigKeyLogLevel) {
	case "debug":
		level = zerolog.DebugLevel
		logLevel = "debug"
	case "info":
		level = zerolog.InfoLevel
		logLevel = "info"
	case "warn":
		level = zerolog.WarnLevel
		logLevel = "warn"
	case "error":
		level = zerolog.ErrorLevel
		logLevel = "error"
	case "fatal":
		level = zerolog.FatalLevel
		logLevel = "fatal"
	case "no":
		level = zerolog.TraceLevel
		logLevel = "no"
	default:
		level = zerolog.InfoLevel
	}
	return level
}

/*
* @Content : 系统级别日志
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-11
 */
func SystemInfo(any ...interface{}) {
	param := localDebug("SystemInfo", "system", any)
	if param == "" {
		param = reqToStr(any)
	}
	nowDate := time.Now().Format(rgconst.GoDateFormat)
	filePath := filePathMerge(logDir, "/"+nowDate, "_SYSTEM.log")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return
	}
	output := zerolog.ConsoleWriter{
		Out:     f,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Local().Format(rgconst.GoTimeFormat)
		},
		FormatLevel: func(i interface{}) string {
			return "| SYSTEMINFO|"
		},
		FormatCaller: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
	}
	logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
	logger.Log().Msg(param)
}

/*
* @Content : 系统级别错误日志
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-11
 */
func SystemError(any ...interface{}) {
	param := localDebug("SystemError", "system", any)
	if param == "" {
		param = reqToStr(any)
	}
	nowDate := time.Now().Format(rgconst.GoDateFormat)
	filePath := filePathMerge(logDir, "/"+nowDate, "_SYSTEM.log")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return
	}
	output := zerolog.ConsoleWriter{
		Out:     f,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Local().Format(rgconst.GoTimeFormat)
		},
		FormatLevel: func(i interface{}) string {
			return "| SYSTEMERROR|"
		},
		FormatCaller: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
	}
	logger := log.Sample(levelSimpler).Output(output).With().Caller().CallerWithSkipFrameCount(3).Logger()
	logger.Log().Msg(param)
}

/*
* @Content : 请求日志
* @Param   :
* @Return  :
* @Author  : LiJunDong
* @Time    : 2022-03-11
 */
func RequestLog(uniqId string, typ string, param string) {
	param = localDebug(typ, uniqId, []interface{}{param})
	nowDate := time.Now().Format(rgconst.GoDateFormat)
	filePath := filePathMerge(logDir, "/"+nowDate, "_REQUEST.log")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return
	}
	output := zerolog.ConsoleWriter{
		Out:     f,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Local().Format(rgconst.GoTimeFormat)
		},
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf("|%s|", typ)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
	}
	logger := log.Sample(levelSimpler).Output(output).With().Logger()
	logger.Log().Fields(map[string]interface{}{"UniqId": uniqId}).Msg(param)
}

func filePathMerge(param ...string) string {
	var build strings.Builder
	for _, item := range param {
		build.WriteString(item)
	}
	result := build.String()
	result = strings.Replace(result, "//", "/", -1)
	return result
}

/*
* @Content : 本地调试日志
* @Param   :
* @Return  : 日志内容
* @Author  : LiJunDong
* @Time    : 2022-03-28
 */
func localDebug(typ string, uniqId string, any []interface{}) (logStr string) {
	if !rgconfig.GetBool(rgconst.ConfigKeyDebug) {
		return ""
	}
	param := reqToStr(any)
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: rgconst.GoTimeFormat,
		NoColor:    true,
		FormatTimestamp: func(i interface{}) string {
			return time.Now().Local().Format(rgconst.GoTimeFormat)
		},
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf("|%s|", typ)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s|", i)
		},
	}
	logger := log.Sample(levelSimpler).Output(output).With().Logger()
	logger.Log().Fields(map[string]interface{}{"UniqId": uniqId}).Msg(param)
	return param
}

func interfaceToString(param interface{}) string {
	thisString := ""
	switch param.(type) {
	case string:
		if value, ok := param.(string); ok {
			thisString = value
		}
	case float64:
		if value, ok := param.(float64); ok {
			thisString = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case int:
		if value, ok := param.(int); ok {
			thisString = strconv.Itoa(value)
		}
	case int64:
		if value, ok := param.(int64); ok {
			thisString = strconv.FormatInt(value, 10)
		}
	case float32:
		if value, ok := param.(float32); ok {
			thisString = strconv.FormatFloat(float64(value), 'f', -1, 32)
		}
	case error:
		if value, ok := param.(error); ok {
			thisString = value.Error()
		}
	case bool:
		if value, ok := param.(bool); ok {
			if value == true {
				thisString = "true"
			} else {
				thisString = "false"
			}
		}
	default:
		logTmp, _ := json.Marshal(param)
		thisString = string(logTmp)
	}
	return thisString
}

func reqToStr(any []interface{}) (data string) {
	if len(any) == 0 {
		return ""
	}
	paramArr := make([]string, 0, 25)
	for key, item := range any {
		if key >= 20 {
			break
		}
		paramArr = append(paramArr, interfaceToString(item))
	}
	data = strings.Join(paramArr, " | ")
	return data
}
