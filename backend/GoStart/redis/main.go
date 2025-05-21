package main

import (
	"context"

	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	traceName = "mxshop-otel"
)

var tp *trace.TracerProvider

func tracerProvider() error {
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("mxshop-user"),
				attribute.String("environment", "dev"),
				attribute.Int("ID", 1),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func main() {
	_ = tracerProvider()
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	cli.AddHook(redisotel.NewTracingHook())
	tr := otel.Tracer("traceName")
	spanCtx, span := tr.Start(context.Background(), "redis")
	cli.Set(spanCtx, "name", "bobby", 0)
	span.End()
	tp.Shutdown(context.Background())
}

/*
nacos一般不会存放配置文件， 这些基本上都会重启服务 - 服务重启会导致系统不可用 k8s YAML
熔断限流 的规则， 1000 user/list, 限流规则改动的会更频繁， 这种改动我们是没有必要去重启服务的
写表， 但是这种做法不太好，这种往往反而放在内容或者redis中更合适，这些数据往往不大
*/
