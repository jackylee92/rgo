package rglog

import (
	"encoding/json"
	"fmt"
	"github.com/jackylee92/rgo/core/rgenv"
	"github.com/jackylee92/rgo/util/rgtime"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
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

type Client struct {
	UniqId string
	sync.Pool
}

type LogLevel string

var (
	LevelDebug    LogLevel = "DEBUG"
	LevelInfo     LogLevel = "INFO"
	LevelWarn     LogLevel = "WARN"
	LevelError    LogLevel = "ERROR"
	LevelRequest  LogLevel = "REQUEST"
	LevelResponse LogLevel = "RESPONSE"
	LevelSystem   LogLevel = "SYSTEM"
	LevelNo       LogLevel = "NO"
)

var (
	fileLoggerClient  zerolog.Logger
	localLoggerClient zerolog.Logger
	logLevel          LogLevel
	logDir            string
)

var fileOutput = zerolog.ConsoleWriter{
	//Out:     getFileOut(),
	NoColor: false,
	FormatTimestamp: func(i interface{}) string {
		return time.Now().Local().Format(rgconst.GoTimeFormat)
	},
	FormatLevel: func(i interface{}) string {
		if i == nil {
			return ""
		}
		if i.(string) == "100" {
			return "| " + string(LevelRequest) + " |"
		}
		if i.(string) == "101" {
			return "| " + string(LevelResponse) + " |"
		}
		if i.(string) == "102" {
			return "| " + string(LevelSystem) + "INFO |"
		}
		if i.(string) == "103" {
			return "| " + string(LevelSystem) + "ERROR |"
		}
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	},
	FormatCaller: func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	},
	FormatMessage: func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	},
	FormatErrFieldName: func(i interface{}) string {
		return "| "
	},
	FormatErrFieldValue: func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	},
}

func Start() {
	SetLogLever(rgconfig.GetStr(rgconst.ConfigKeyLogLevel))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
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
	sampler := &zerolog.BurstSampler{
		Burst:       5000,
		Period:      time.Second,
		NextSampler: &zerolog.BasicSampler{N: 100},
	}
	logSampled := &zerolog.LevelSampler{
		DebugSampler: sampler,
		TraceSampler: sampler,
		InfoSampler:  sampler,
		WarnSampler:  sampler,
		ErrorSampler: sampler,
	}
	fileOutput.Out = getFileOut()
	fileLoggerClient = zerolog.New(fileOutput).Sample(logSampled).With().Timestamp().Stack().CallerWithSkipFrameCount(3).Logger()
	localLoggerClient = zerolog.New(os.Stdout).Sample(logSampled).With().Timestamp().Stack().CallerWithSkipFrameCount(4).Logger()
	backUpLog()
}

func New(uniqId string) *Client {
	return &Client{UniqId: uniqId}
}

func (c *Client) Info(any ...interface{}) {
	param := localDebug(LevelInfo, c.UniqId, any)
	if !log.Info().Enabled() {
		return
	}
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.Info().Msg(c.UniqId + " | " + param)
}

func (c *Client) Error(any ...interface{}) {
	param := localDebug(LevelError, c.UniqId, any)
	if !log.Error().Enabled() {
		return
	}
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.Error().Err(errors.New(param)).Msg(c.UniqId)
}

func (c *Client) Debug(any ...interface{}) {
	param := localDebug(LevelDebug, c.UniqId, any)
	if !log.Debug().Enabled() {
		return
	}
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.Debug().Msg(c.UniqId + " | " + param)
}

func (c *Client) Warn(any ...interface{}) {
	param := localDebug(LevelWarn, c.UniqId, any)
	if !log.Warn().Enabled() {
		return
	}
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.Warn().Msg(c.UniqId + " | " + param)
}

func SetLogLever(param string) {
	switch param {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logLevel = LevelDebug
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logLevel = LevelInfo
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		logLevel = LevelWarn
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logLevel = LevelError
	case "no":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		logLevel = LevelNo
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logLevel = LevelInfo
	}
	return
}

func SystemInfo(any ...interface{}) {
	param := localDebug(LevelSystem, "SystemInfo", any)
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.WithLevel(102).Msg(param)
}

func SystemError(any ...interface{}) {
	param := localDebug(LevelSystem, "SystemError", any)
	if param == "" {
		param = reqToStr(any)
	}
	fileLoggerClient.WithLevel(103).Msg(param)
}

func RequestLog(uniqId string, typ LogLevel, any ...interface{}) {
	param := localDebug(typ, uniqId, any)
	if param == "" {
		param = reqToStr(any)
	}
	var requestLevel zerolog.Level = 100
	if typ == LevelResponse {
		requestLevel = 101
	}
	fileLoggerClient.WithLevel(requestLevel).Msg(uniqId + " | " + param)
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

func localDebug(typ LogLevel, uniqId string, any []interface{}) (logStr string) {
	if !rgconfig.GetBool(rgconst.ConfigKeyDebug) {
		return ""
	}
	param := reqToStr(any)
	switch typ {
	case LevelDebug:
		localLoggerClient.Debug().Msg(uniqId + " | " + param)
	case LevelInfo:
		localLoggerClient.Info().Msg(uniqId + " | " + param)
	case LevelWarn:
		localLoggerClient.Warn().Msg(uniqId + " | " + param)
	case LevelError:
		localLoggerClient.Error().Err(errors.New(param)).Msg(uniqId)
		return
	case LevelRequest:
		localLoggerClient.WithLevel(100).Msg(uniqId + " | " + param)
	case LevelResponse:
		localLoggerClient.WithLevel(101).Msg(uniqId + " | " + param)
	case LevelSystem:
		var newLevel zerolog.Level = 102
		if typ == "SystemInfo" {
			newLevel = 103
		}
		localLoggerClient.WithLevel(newLevel).Msg(param)
	}
	return
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

func GetLogLevel() string {
	return string(logLevel)
}

func getFileOut() (f *os.File) {
	f, err := os.OpenFile(logDir+"/"+rgenv.GetAppName()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("日志文件打开失败:", err)
	}
	return f
}

func backUpLog() {
	backJob := cron.New()
	spec := "@daily"
	//spec := "@every 1m"
	err := backJob.AddFunc(spec, func() {
		copyCleanLogFile()
	})
	if err != nil {
		fmt.Println("日志文件备份任务启动失败:", err)
	}
	backJob.Start()
}

func copyCleanLogFile() {
	//filePath := logDir + "/" + rgenv.GetAppName() + ".log." + time.Now().Format("2006-01-02 15:04")
	filePath := logDir + "/" + rgenv.GetAppName() + ".log." + rgtime.NowDate()
	err := os.Rename(logDir+"/"+rgenv.GetAppName()+".log", filePath)
	if err != nil {
		fmt.Println("日志文件重命名失败:", err)
	}
	f, err := os.OpenFile(logDir+"/"+rgenv.GetAppName()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("日志文件打开失败:", err)
	}
	if fileOutput.Out != nil {
		if of, ok := fileOutput.Out.(*os.File); ok {
			_ = of.Close()
		}
	}
	fileOutput.Out = f
	fileLoggerClient = fileLoggerClient.Output(fileOutput)
}
