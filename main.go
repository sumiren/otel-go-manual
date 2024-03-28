package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var err error
	if err != nil {
		log.Fatal(err)
	}

	otelShutdown, err, tracerProvider := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-northeast-1.prod.aws.tidbcloud.com",
	})

	r := setupRouter(tracerProvider)
	// Listen and Server in 0.0.0.0:8080
	fmt.Println("start")
	go func() {
		r.Run(":8080")
	}()
	<-ctx.Done()

	fmt.Println("end")

}

func setupRouter(tracerProvider *trace.TracerProvider) *gin.Engine {

	// Disable Console Color
	// gin.DisableConsoleColor()

	r := gin.Default()
	r.Use(otelgin.Middleware("otel-go"))

	db, _ := otelsql.Open("mysql", "3uAcMNCBzHP6Q4A.root:WFCUi393BPAf1vmN@tcp(gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000)/main?tls=tidb", otelsql.WithAttributes(
		semconv.DBSystemMySQL))

	// Ping test
	r.GET("/sql", func(c *gin.Context) {
		results := dbquery(c, db)
		anotherResult, _ := httpget(c)

		c.JSON(http.StatusOK, map[string]interface{}{
			"db":   results,
			"http": anotherResult,
		})
	})

	return r
}
