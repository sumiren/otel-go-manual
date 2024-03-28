package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func httpget(c *gin.Context) (map[string]interface{}, error) {
	//_, span := tracer.Start(c.Request.Context(), "dbquery", trace.WithAttributes(attribute.String("id", "hogehoge")))
	//defer span.End()

	tracer := otel.GetTracerProvider()
	_, span := tracer.Tracer("application").Start(c.Request.Context(), "custom httpget method", trace.WithAttributes(attribute.String("id", "hogehoge")))
	defer span.End()

	// HTTP GETリクエストを送信
	resp, _ := otelhttp.Get(c.Request.Context(), "https://www.yahoo.co.jp")
	defer resp.Body.Close()

	result := map[string]interface{}{
		"result": resp.StatusCode,
	}

	return result, nil
}
