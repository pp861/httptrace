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

// http client do with trace
// start a new span only if there is a parent span in context.
func DoHttpSend(ctx context.Context, client *http.Client, req *http.Request) (rsp *http.Response, err error) {
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentContext := parent.Context()
		span := opentracing.StartSpan(
			"HttpClient Call "+req.URL.RequestURI(),
			opentracing.ChildOf(parentContext.(opentracing.SpanContext)),
			opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			ext.SpanKindRPCClient,
		)

		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL.String())

		defer span.Finish()

		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			span.LogFields(opentracingLog.String("inject-error", err.Error()))
		}
	}

	return client.Do(req)
}

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
