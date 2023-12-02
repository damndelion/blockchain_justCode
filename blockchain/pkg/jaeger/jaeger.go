package jaeger

import (
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

func InitJaeger(serviceName string, url string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 100, // 100 traces per second
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: url,
		},
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.New(serviceName)

	return Tracer, closer, err
}

func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request, funcDesc string) opentracing.Span {
	spanCtx, err := Extract(tracer, r)
	if err != nil {
		return nil
	}

	return tracer.StartSpan(funcDesc, ext.RPCServerOption(spanCtx))
}

func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}
