package httptrace

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// init trace
// service: server name
// exporterAddr: trace agent addr("192.168.1.10:6831")
func InitTrace(service, exporterAddr string) (opentracing.Tracer, io.Closer, error) {
	if exporterAddr == "" {
		return nil, nil, errors.New("trace exporterAddr empty")
	}

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: exporterAddr,
		},
	}

	return cfg.New(service, config.Logger(jaeger.StdLogger))
}

// http client do with trace text map
// start a new span only if there is a parent span in context.
// func DoHttpSendWithTextMap(textMap map[string]string, client *http.Client, req *http.Request) (rsp *http.Response, err error) {
// 	parentContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(textMap))
// 	if parentContext == nil {
// 		return client.Do(req)
// 	}

// 	return httpSend(parentContext, client, req)
// }

// http client do with trace
// start a new span only if there is a parent span in context.
func DoHttpSend(ctx context.Context, client *http.Client, req *http.Request) (rsp *http.Response, err error) {
	if client == nil || req == nil {
		return nil, errors.New("httptrace: httpClient or httpReq can not be nil")
	}

	if ctx == nil {
		return client.Do(req)
	}

	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		return client.Do(req)
	}

	return httpSend(parent.Context(), client, req)
}

func httpSend(parentContext opentracing.SpanContext, client *http.Client, req *http.Request) (rsp *http.Response, err error) {
	if parentContext == nil {
		return nil, errors.New("httptrace: SpanContext is nil")
	}

	span := opentracing.StartSpan(
		"HttpClient Call "+req.URL.RequestURI(),
		opentracing.ChildOf(parentContext.(opentracing.SpanContext)),
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCClient,
	)

	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())

	defer span.Finish()

	err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		span.LogFields(opentracingLog.String("inject-error", err.Error()))
	}

	rsp, err = client.Do(req)
	if err != nil {
		ext.LogError(span, err)
	}

	return rsp, err
}
