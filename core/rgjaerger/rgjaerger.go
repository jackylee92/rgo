package rgjaerger

import (
	"errors"
	"go.opentelemetry.io/otel/bridge/opentracing"
	"io"
	"net/http"
	"time"

	"rgo/core/rgconfig"
	"rgo/core/rgglobal"
	"rgo/core/rgglobal/rgconst"
	"rgo/core/rgglobal/rgerror"
	"rgo/core/rglog"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	// jaegerlog "github.com/uber/jaeger-client-go/log"
)

// 是否开启
var jaergerStatus bool

func JaergerStatus() bool {
	return jaergerStatus
}

func SetJaergerStatus(status bool) {
	if status {
		rglog.SystemInfo("启动项【jaerger_status】:开启")
	}
	jaergerStatus = status
}

type Client struct {
	ctx    opentracing.SpanContext
	tracer opentracing.Tracer
}

// 从上下文中解析并创建一个新的 trace，获得传播的 上下文(SpanContext)
func GetTracer(header http.Header) (opentracing.Tracer, opentracing.SpanContext, io.Closer, error) {
	tracer, closer, err := createTracer()
	if err != nil {
		return nil, nil, nil, err
	}
	// 继承别的进程传递过来的上下文
	spanContext, _ := tracer.Extract(opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header))
	return tracer, spanContext, closer, err
}

// 创建一个Tracer
func createTracer() (opentracing.Tracer, io.Closer, error) {
	appName := rgglobal.AppName
	jaergerHost := rgconfig.GetStr(rgconst.ConfigKeyJaergerHost)
	if jaergerHost == "" {
		return nil, nil, errors.New(rgerror.ErrorJaergerHostNil)
	}
	if appName == "" {
		return nil, nil, errors.New(rgerror.ErrorAppNameNil)
	}
	var cfg = jaegercfg.Configuration{
		ServiceName: appName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaergerHost,
		},
	}

	// jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		// jaegercfg.Logger(jLogger),
		jaegercfg.MaxTagValueLength(65535),
	)
	return tracer, closer, err
}

func New(ctx opentracing.SpanContext, tracer opentracing.Tracer) *Client {
	if !jaergerStatus || ctx == nil || tracer == nil {
		return nil
	} else {
		return &Client{
			ctx:    ctx,
			tracer: tracer,
		}
	}
}

func (c *Client) Send(action string, param map[string]interface{}, parentSpanInterface interface{}) opentracing.SpanContext {
	if !jaergerStatus {
		return nil
	}
	param["time"] = time.Now().Format(rgconst.GoTimeFormat)
	param["time.nano"] = time.Now().UnixNano()
	param["local.ip"] = rgglobal.LocalIp
	if len(param) == 0 {
		return nil
	}
	tags := opentracing.Tags{}
	for key, value := range param {
		tags[key] = value
	}
	span := opentracing.ChildOf(c.ctx)
	parentSpan, ok := parentSpanInterface.(opentracing.SpanContext)
	if parentSpan != nil && ok {
		span = opentracing.ChildOf(parentSpan)
	}
	childSpan := c.tracer.StartSpan(action, span, tags)
	defer childSpan.Finish() // 可手动调用 Finish()
	return childSpan.Context()
}
